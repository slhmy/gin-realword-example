package core

import (
	"context"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lmittmann/tint"
)

const (
	logLevelConfigKey  = "log.level"
	logFormatConfigKey = "log.format"

	slogFields ContextKey = "slog_fields"
)

var (
	logLevel  slog.Level
	logFormat string
)

func stringToSlogLevel(
	level string,
) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type contextHandler struct {
	slog.Handler
}

func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	if len(ServiceRole) > 0 {
		r.AddAttrs(slog.String("service_role", ServiceRole))
	}
	if len(hostname) > 0 {
		r.AddAttrs(slog.String("hostname", hostname))
	}
	if attrs, ok := ctx.Value(string(slogFields)).([]slog.Attr); ok {
		for _, v := range attrs {
			r.AddAttrs(v)
		}
	}
	return h.Handler.Handle(ctx, r)
}

func AppendLogFieldToCtx(parent context.Context, attr slog.Attr) context.Context {
	if v, ok := parent.Value(slogFields).([]slog.Attr); ok {
		v = append(v, attr)
		return context.WithValue(parent, slogFields, v)
	}
	v := []slog.Attr{}
	v = append(v, attr)
	return context.WithValue(parent, slogFields, v)
}

func AppendLogFieldToGinCtx(
	parent *gin.Context, attr slog.Attr,
) {
	v, ok := parent.Value(slogFields).([]slog.Attr)
	if !ok {
		v = []slog.Attr{}
	}
	v = append(v, attr)
	parent.Set(string(slogFields), v)
}

func init() {
	logLevel = stringToSlogLevel(ConfigStore.GetString(logLevelConfigKey))
	logFormat = ConfigStore.GetString(logFormatConfigKey)

	var handler slog.Handler
	switch logFormat {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	case "plain":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	default:
		handler = tint.NewHandler(os.Stdout, &tint.Options{Level: logLevel})
	}
	handler = &contextHandler{handler}
	slog.SetDefault(slog.New(handler))
}
