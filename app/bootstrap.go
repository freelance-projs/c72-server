package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/ngoctd314/c72-api-server/app/route"
	"github.com/ngoctd314/c72-api-server/pkg/helper"
	"github.com/ngoctd314/c72-api-server/pkg/repository"
	"github.com/ngoctd314/common/env"
	"github.com/ngoctd314/common/net/conn"
	"github.com/ngoctd314/common/net/ghttp"
	"gorm.io/gorm"
)

type app struct {
	httpServer *ghttp.Server
	dbConn     *gorm.DB
}

func New(ctx context.Context) *app {
	initLogger()

	// auto migrate up
	if err := migrateUp(); err != nil {
		dsn := env.GetString("mysql.tag_scan.dsn")
		slog.Error("error occur when migrate up", "err", err, "dsn", dsn)
	}

	dbConn := mustInitDBConn()
	tagRepo := repository.NewTag(dbConn)

	handler := route.Handler(tagRepo)

	return &app{
		httpServer: mustInitServer(ctx, handler),
		dbConn:     dbConn,
	}
}

func (i *app) Start(ctx context.Context) {
	localIP, _ := helper.LocalIP()
	slog.Info("Server is starting", "local_ip", localIP)
	if err := i.httpServer.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			slog.Info("ErrServerClosed")
			return
		}
		slog.Error("error occur when ListenAndServe", "err", err)
		panic(err)
	}
}

func (i *app) Shutdown(ctx context.Context) error {
	if err := i.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("error occur when shutdown http server, %w", err)
	}

	dbConn, err := i.dbConn.DB()
	if err != nil {
		return fmt.Errorf("error occur when get db connection, %w", err)
	}
	if err := dbConn.Close(); err != nil {
		return fmt.Errorf("error occur when close db connection, %w", err)
	}

	return nil
}

func mustInitServer(ctx context.Context, usecase http.Handler) *ghttp.Server {
	server, err := ghttp.NewServer(usecase,
		// inject main ctx
		ghttp.WithBaseContext(func(_ net.Listener) context.Context {
			return ctx
		}),
		ghttp.WithServerLogger(slog.Default()),
	)
	if err != nil {
		panic(err)
	}

	return server
}

func mustInitDBConn() *gorm.DB {
	sqlConn, err := conn.SQL(conn.MySQLDriver, "tag_scan")
	if err != nil {
		panic(err)
	}

	db, err := repository.DB(sqlConn)
	if err != nil {
		panic(err)
	}

	return db
}
