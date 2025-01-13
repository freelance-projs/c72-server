package app

import (
	"io"
	"log/slog"
	"os"

	"github.com/ngoctd314/common/env"
)

func initLogger() {
	replaceAttr := func(groups []string, a slog.Attr) slog.Attr {
		return a
	}

	opts := &slog.HandlerOptions{
		AddSource:   true,
		Level:       slog.LevelInfo,
		ReplaceAttr: replaceAttr,
	}

	slogHandler := loggerHandler(loggerWriter(), opts)
	slog.SetDefault(slog.New(slogHandler))
}

func loggerHandler(writer io.Writer, opts *slog.HandlerOptions) slog.Handler {
	switch env.GetString("logger.format") {
	case "json":
		return slog.NewJSONHandler(writer, opts)
	case "text":
		return slog.NewTextHandler(writer, opts)
	default:
		return slog.NewTextHandler(writer, opts)
	}
}

func loggerWriter() io.Writer {
	return os.Stdout
}
