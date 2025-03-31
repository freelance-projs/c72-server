package usecase

import (
	"context"
	"log/slog"

	"github.com/ngoctd314/common/net/ghttp"
)

type ping struct {
}

func Ping() *ping {
	return &ping{}
}

// Usecase implements ghttp.Usecase.
func (p *ping) Usecase(ctx context.Context, req *struct{}) (*ghttp.ResponseBody, error) {
	slog.Info("ping")
	return &ghttp.ResponseBody{
		Success: true,
		Data:    req,
	}, nil
}
