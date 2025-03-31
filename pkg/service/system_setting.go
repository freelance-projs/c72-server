package service

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"github.com/ngoctd314/c72-api-server/pkg/model"
)

type settingRepo interface {
	ListSetting(ctx context.Context) ([]model.Setting, error)
}

type systemSetting struct {
	settingRepo settingRepo
	data        map[string]string
	rwm         sync.RWMutex
}

func newSystemSetting(settingRepo settingRepo) *systemSetting {
	ss := &systemSetting{
		settingRepo: settingRepo,
	}
	ss.Refresh(context.Background())

	go func() {
		refreshDuration := 5 * time.Minute
		ticker := time.NewTicker(refreshDuration)
		for range ticker.C {
			if err := ss.Refresh(context.Background()); err != nil {
				slog.Error("error refreshing system settings", "err", err)
				return
			}
			ticker.Reset(refreshDuration)
		}
	}()

	return ss
}

var SystemSetting *systemSetting

func InitSystemSetting(settingRepo settingRepo) {
	ss := newSystemSetting(settingRepo)
	SystemSetting = ss
}

var errSettingNotFound = errors.New("setting not found")

const (
	SettingKeyTXLockTime = "tx_lock_time"
)

func (s *systemSetting) GetTxLockTime() time.Duration {
	s.rwm.RLock()
	defer s.rwm.RUnlock()

	v, ok := s.data[SettingKeyTXLockTime]
	if ok {
		d, err := strconv.Atoi(v)
		if err != nil {
			return 0
		}

		return time.Duration(d) * time.Minute
	}
	return 0
}

func (s *systemSetting) Refresh(ctx context.Context) error {
	s.rwm.Lock()
	defer s.rwm.Unlock()

	settings, err := s.settingRepo.ListSetting(ctx)
	if err != nil {
		return err
	}

	s.data = make(map[string]string)
	for _, setting := range settings {
		s.data[setting.Key] = setting.Value
	}

	return nil
}
