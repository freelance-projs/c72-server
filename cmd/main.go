package main

import (
	"context"
	"flag"
	"log/slog"

	"github.com/ngoctd314/c72-api-server/app"
	"github.com/ngoctd314/common/core"
	"github.com/ngoctd314/common/env"
)

var (
	configFile = flag.String("config", "config.yaml", "Path to configuration file.")
)

func main() {
	flag.Parse()

	env.Init(env.WithFile(*configFile))

	ctx := context.Background()
	app := app.New(ctx)

	instance := core.NewInstance(ctx, app, core.WithLogger(slog.Default()))
	instance.Bootstrap()
}
