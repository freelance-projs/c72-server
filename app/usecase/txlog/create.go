package txlog

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/c72-api-server/pkg/repository"
	"github.com/ngoctd314/common/env"
	"github.com/ngoctd314/common/net/ghttp"
	"gopkg.in/telebot.v4"
)

type create struct {
	sheetSvc *sheetService
	repo     *repository.Repository
	teleBot  *telebot.Bot
	sheetID  string
}

func Create(repo *repository.Repository) *create {
	pref := telebot.Settings{
		Token: env.GetString("telegram.botToken"),
	}
	teleBot, err := telebot.NewBot(pref)
	if err != nil {
		panic(err)
	}
	sheetSvc := newSheetService()

	uc := &create{
		repo:     repo,
		teleBot:  teleBot,
		sheetSvc: sheetSvc,
	}
	setting, err := repo.GetSetting(context.Background())
	if err == nil {
		uc.sheetID = setting.TxLogSheetID
	} else {
		uc.sheetID = "1XIMAojHp1g-SMt8aOY-IGaPB0hT-KnW1HsBWB4VcV64"
	}

	go func() {
		ticker := time.NewTicker(time.Minute)
		for range ticker.C {
			setting, err := repo.GetSetting(context.Background())
			slog.Info("updating tx log sheet id")
			if err == nil {
				uc.sheetID = setting.TxLogSheetID
			}
			ticker.Reset(time.Minute)
		}
	}()

	return uc

}

func (uc *create) Usecase(ctx context.Context, req *dto.CreateTxLogRequest) (*ghttp.ResponseBody, error) {
	switch req.Action {
	case "lending":
		if err := uc.createLendingTx(ctx, req.Entity, req.TagIDs); err != nil {
			return nil, err
		}
		// send telegram

		return ghttp.ResponseBodyOK(nil, ghttp.ResponseBodyWithMessage(
			fmt.Sprintf("%s mượn thành công %d vật phẩm", req.Entity, len(req.TagIDs)),
		)), nil

	case "lending_return":
		if err := uc.createLendingReturnTx(ctx, req.Entity, req.TagIDs); err != nil {
			return nil, err
		}
		return ghttp.ResponseBodyOK(nil, ghttp.ResponseBodyWithMessage(
			fmt.Sprintf("%s trả thành công %d vật phẩm", req.Entity, len(req.TagIDs)),
		)), nil

	case "washing":
		if err := uc.createWashingTx(ctx, req.Entity, req.TagIDs); err != nil {
			return nil, err
		}
		return ghttp.ResponseBodyOK(nil, ghttp.ResponseBodyWithMessage(
			fmt.Sprintf("Giao %d vật phẩm cho %s thành công", len(req.TagIDs), req.Entity),
		)), nil

	case "washing_return":
		if err := uc.createWashingReturnTx(ctx, req.Entity, req.TagIDs); err != nil {
			return nil, err
		}
		return ghttp.ResponseBodyOK(nil, ghttp.ResponseBodyWithMessage(
			fmt.Sprintf("Nhận %d vật phẩm từ %s thành công", len(req.TagIDs), req.Entity),
		)), nil
	}

	return nil, nil
}

type lendingTxSheet struct {
	Department string `header:"Phòng ban"`
	Action     string `header:"Hàng động"`
	TagName    string `header:"Tên vật phẩm"`
	Count      int    `header:"Số lượng"`
	CreateAt   string `header:"Vào lúc"`
}

func (uc *create) createLendingTx(ctx context.Context, department string, tagIDs []string) error {
	txID, txLogDepartment, err := uc.repo.CreateLendingTx(ctx, department, tagIDs)
	if err != nil {
		return err
	}

	sheetCols := make([]any, 0, len(txLogDepartment))
	for _, v := range txLogDepartment {
		action := "Mượn"
		if v.Action == model.EDepartmentActionReturned {
			action = "Trả"
		}
		sheetCols = append(sheetCols, lendingTxSheet{
			Department: v.Department,
			Action:     action,
			TagName:    v.TagName,
			Count:      int(v.Lending),
			CreateAt:   v.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	if len(sheetCols) > 0 {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("error inserting data to sheet", "err", r)
				}
			}()
			now := time.Now()
			sheetName := "Nội bộ " + time.Now().Format("2006-01-02")
			if err := uc.sheetSvc.insert(uc.sheetID, sheetName, sheetCols); err != nil {
				slog.Error("error inserting data to sheet", "err", err)
				return
			}
			slog.Info("insert data to sheet successfully", "sheetID", uc.sheetID, "since", time.Since(now).Seconds())
		}()
	}

	go func() {
		if _, err := uc.teleBot.Send(&telebot.Chat{
			ID: -1002500429787,
		}, fmt.Sprintf(`Ghi nhận giao dịch Khoa Mượn Đồ
		- Thực hiện lúc %s
		- %s mượn %d vật phẩm

	Xem chi tiết tại:
		- %s`,
			time.Now().Format("15:04:05"),
			department,
			len(tagIDs),
			fmt.Sprintf("http://154.26.134.232:5081/tx-log/department/%d", txID)),
			&telebot.SendOptions{},
		); err != nil {
			slog.Error("error sending message", "err", err)
		}
	}()

	return nil
}

func (uc *create) createLendingReturnTx(ctx context.Context, department string, tagIDs []string) error {
	txIDs, txLogDepartment, err := uc.repo.ReturnLendingTx(ctx, department, tagIDs)
	if err != nil {
		return err
	}

	sheetCols := make([]any, 0, len(txLogDepartment))
	for _, v := range txLogDepartment {
		action := "Mượn"
		if v.Action == model.EDepartmentActionReturned {
			action = "Trả"
		}
		sheetCols = append(sheetCols, lendingTxSheet{
			Department: v.Department,
			Action:     action,
			TagName:    v.TagName,
			Count:      int(v.Returned),
			CreateAt:   v.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	if len(sheetCols) > 0 {
		go func() {
			now := time.Now()
			sheetName := "Nội bộ " + time.Now().Format("2006-01-02")
			if err := uc.sheetSvc.insert(uc.sheetID, sheetName, sheetCols); err != nil {
				slog.Error("error inserting data to sheet", "err", err)
				return
			}
			slog.Info("insert data to sheet successfully", "sheetID", uc.sheetID, "since", time.Since(now).Seconds())
		}()
	}

	var links []string
	for _, txID := range txIDs {
		link := fmt.Sprintf("http://154.26.134.232:5081/tx-log/department/%d", txID)
		links = append(links, link)
	}
	go func() {
		if _, err := uc.teleBot.Send(&telebot.Chat{
			ID: -1002500429787,
		}, fmt.Sprintf(`Ghi nhận giao dịch Khoa Trả Đồ
	- Thực hiện lúc %s
	- %s trả %d vật phẩm

Xem chi tiết tại:
	- %s`,
			time.Now().Format("15:04:05"),
			department,
			len(tagIDs),
			strings.Join(links, "\n\t- "),
		),
			&telebot.SendOptions{},
		); err != nil {
			slog.Error("error sending message", "err", err)
		}
	}()

	return nil
}

type washingTxSheet struct {
	Company  string `header:"Công ty"`
	Action   string `header:"Hàng động"`
	TagName  string `header:"Tên vật phẩm"`
	Count    int    `header:"Số lượng"`
	CreateAt string `header:"Vào lúc"`
}

func (uc *create) createWashingTx(ctx context.Context, department string, tagIDs []string) error {
	txID, txLogCompany, err := uc.repo.CreateWashingTx(ctx, department, tagIDs)
	if err != nil {
		return err
	}

	sheetCols := make([]any, 0, len(txLogCompany))
	for _, v := range txLogCompany {
		action := "Giặt"
		if v.Action == model.ECompanyActionReturned {
			action = "Trả"
		}
		sheetCols = append(sheetCols, washingTxSheet{
			Company:  v.Company,
			Action:   action,
			TagName:  v.TagName,
			Count:    int(v.Washing),
			CreateAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	if len(sheetCols) > 0 {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("error inserting data to sheet", "err", r)
				}
			}()
			now := time.Now()
			sheetName := "Công ty " + time.Now().Format("2006-01-02")
			if err := uc.sheetSvc.insert(uc.sheetID, sheetName, sheetCols); err != nil {
				slog.Error("error inserting data to sheet", "err", err)
				return
			}
			slog.Info("công ty giặt đồ thành công", "sheetID", uc.sheetID, "since", time.Since(now).Seconds())
		}()
	}

	go func() {
		if _, err := uc.teleBot.Send(&telebot.Chat{
			ID: -1002500429787,
		}, fmt.Sprintf(`Ghi nhận giao dịch Giặt Đồ Bẩn 
	- Thực hiện lúc %s
	- %s giặt %d vật phẩm

Xem chi tiết tại:
	- %s`,
			time.Now().Format("15:04:05"),
			department,
			len(tagIDs),
			fmt.Sprintf("http://154.26.134.232:5081/tx-log/company/%d", txID)),
			&telebot.SendOptions{},
		); err != nil {
			slog.Error("error sending message", "err", err)
		}
	}()

	return nil
}

func (uc *create) createWashingReturnTx(ctx context.Context, department string, tagIDs []string) error {
	txIDs, txLogCompany, err := uc.repo.ReturnWashingTx(ctx, department, tagIDs)
	if err != nil {
		return err
	}

	sheetCols := make([]any, 0, len(txLogCompany))
	for _, v := range txLogCompany {
		action := "Giặt"
		if v.Action == model.ECompanyActionReturned {
			action = "Trả"
		}
		sheetCols = append(sheetCols, washingTxSheet{
			Company:  v.Company,
			Action:   action,
			TagName:  v.TagName,
			Count:    int(v.Returned),
			CreateAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	if len(sheetCols) > 0 {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("error inserting data to sheet", "err", r)
				}
			}()
			now := time.Now()
			sheetName := "Công ty " + time.Now().Format("2006-01-02")
			if err := uc.sheetSvc.insert(uc.sheetID, sheetName, sheetCols); err != nil {
				slog.Error("error inserting data to sheet", "err", err)
				return
			}
			slog.Info("công ty trả đồ thành công", "sheetID", uc.sheetID, "since", time.Since(now).Seconds())
		}()
	}

	var links []string
	for _, txID := range txIDs {
		link := fmt.Sprintf("http://154.26.134.232:5081/tx-log/company/%d", txID)
		links = append(links, link)
	}
	go func() {
		if _, err := uc.teleBot.Send(&telebot.Chat{
			ID: -1002500429787,
		}, fmt.Sprintf(`Ghi nhận giao dịch Trả Đồ Sạch
	- Thực hiện lúc %s
	- %s trả %d vật phẩm

Xem chi tiết tại:
	- %s`,
			time.Now().Format("15:04:05"),
			department,
			len(tagIDs),
			strings.Join(links, "\n\t- "),
		),
			&telebot.SendOptions{},
		); err != nil {
			slog.Error("error sending message", "err", err)
		}
	}()

	return nil
}
