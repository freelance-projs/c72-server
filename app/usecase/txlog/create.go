package txlog

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/repository"
	"github.com/ngoctd314/common/env"
	"github.com/ngoctd314/common/net/ghttp"
	"gopkg.in/telebot.v4"
)

type create struct {
	repo    *repository.Laundry
	teleBot *telebot.Bot
}

func Create(repo *repository.Laundry) *create {
	pref := telebot.Settings{
		Token: env.GetString("telegram.botToken"),
	}
	teleBot, err := telebot.NewBot(pref)
	if err != nil {
		panic(err)
	}

	return &create{
		repo:    repo,
		teleBot: teleBot,
	}
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

func (uc *create) createLendingTx(ctx context.Context, department string, tagIDs []string) error {
	mTxLog, err := uc.repo.CreateLendingTx(ctx, department, tagIDs)
	if err != nil {
		return err
	}

	txID := mTxLog.ID
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
			fmt.Sprintf("http://154.26.134.232:5081/lending/%d", txID)),
			&telebot.SendOptions{},
		); err != nil {
			slog.Error("error sending message", "err", err)
		}
	}()

	return nil
}

func (uc *create) createLendingReturnTx(ctx context.Context, department string, tagIDs []string) error {
	txIDs, err := uc.repo.ReturnLendingTx(ctx, department, tagIDs)
	if err != nil {
		return err
	}

	var links []string
	for _, txID := range txIDs {
		link := fmt.Sprintf("http://154.26.134.232:5081/lending/%d", txID)
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

func (uc *create) createWashingTx(ctx context.Context, department string, tagIDs []string) error {
	mTxLog, err := uc.repo.CreateWashingTx(ctx, department, tagIDs)
	if err != nil {
		return err
	}

	txID := mTxLog.ID
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
			fmt.Sprintf("http://154.26.134.232:5081/washing/%d", txID)),
			&telebot.SendOptions{},
		); err != nil {
			slog.Error("error sending message", "err", err)
		}
	}()

	return nil
}

func (uc *create) createWashingReturnTx(ctx context.Context, department string, tagIDs []string) error {
	txIDs, err := uc.repo.ReturnWashingTx(ctx, department, tagIDs)
	if err != nil {
		return err
	}

	var links []string
	for _, txID := range txIDs {
		link := fmt.Sprintf("http://154.26.134.232:5081/washing/%d", txID)
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
