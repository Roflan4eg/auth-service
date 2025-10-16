package logger

import (
	"context"
	"errors"
	"github.com/Roflan4eg/auth-serivce/config"
	"log/slog"
	"os"
	"time"
)

type ctxKey string

const (
	slogFields ctxKey = "slog_fields"
)

type logCtx struct {
	traceID string
	userID  string
	method  string
	data    map[string]any
}

type Logger struct {
	*slog.Logger
}

func NewLogger(cfg *config.AppConfig) *Logger {
	var handler slog.Handler

	switch cfg.Environment {
	case "development": //change writing dir
		handler = NewHandlerMiddleware(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		}))
	case "local":
		handler = NewHandlerMiddleware(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		}))

	case "production":
		handler = NewHandlerMiddleware(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: true,
		}))
	}

	logger := slog.New(handler)

	logger = logger.With(slog.Group(
		"app",
		slog.String("name", cfg.Name),
		slog.String("version", cfg.Version),
	))

	return &Logger{logger}

}

func (l *Logger) String(key, value string) slog.Attr {
	return slog.String(key, value)
}

func (l *Logger) Bool(key string, value bool) slog.Attr {
	return slog.Bool(key, value)
}

func (l *Logger) Int(key string, value int) slog.Attr {
	return slog.Int(key, value)
}

func (l *Logger) Uint64(key string, value uint) slog.Attr {
	return slog.Uint64(key, uint64(value))
}

func (l *Logger) Any(key string, value any) slog.Attr {
	return slog.Any(key, value)
}

func (l *Logger) Duration(key string, value time.Duration) slog.Attr {
	return slog.Duration(key, value)
}

type HandlerMiddleware struct {
	next slog.Handler
}

func NewHandlerMiddleware(next slog.Handler) *HandlerMiddleware {
	return &HandlerMiddleware{next: next}
}

func (h *HandlerMiddleware) Enabled(ctx context.Context, rec slog.Level) bool {
	return h.next.Enabled(ctx, rec)
}

func (h *HandlerMiddleware) Handle(ctx context.Context, rec slog.Record) error {
	if c, ok := ctx.Value(slogFields).(logCtx); ok {
		if c.traceID != "" {
			rec.Add("traceID", c.traceID)
		}
		if c.userID != "" {
			rec.Add("userID", c.userID)
		}
		if c.method != "" {
			rec.Add("method", c.method)
		}
		if c.data != nil {
			rec.Add("data", c.data)
		}
	}
	return h.next.Handle(ctx, rec)
}

func (h *HandlerMiddleware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &HandlerMiddleware{next: h.next.WithAttrs(attrs)} // обернуть
}

func (h *HandlerMiddleware) WithGroup(name string) slog.Handler {
	return &HandlerMiddleware{next: h.next.WithGroup(name)} // обернуть
}

func WithUserID(ctx context.Context, userID string) context.Context {
	if c, ok := ctx.Value(slogFields).(logCtx); ok {
		c.userID = userID
		return context.WithValue(ctx, slogFields, c)
	}
	return context.WithValue(ctx, slogFields, logCtx{userID: userID})
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	if c, ok := ctx.Value(slogFields).(logCtx); ok {
		c.traceID = traceID
		return context.WithValue(ctx, slogFields, c)
	}
	return context.WithValue(ctx, slogFields, logCtx{traceID: traceID})
}

func WithMethod(ctx context.Context, method string) context.Context {
	if c, ok := ctx.Value(slogFields).(logCtx); ok {
		c.method = method
		return context.WithValue(ctx, slogFields, c)
	}
	return context.WithValue(ctx, slogFields, logCtx{method: method})
}

func WithData(ctx context.Context, data map[string]any) context.Context {
	if c, ok := ctx.Value(slogFields).(logCtx); ok {
		c.data = maskSensitiveData(data)
		return context.WithValue(ctx, slogFields, c)
	}
	return context.WithValue(ctx, slogFields, logCtx{data: data})
}

func maskSensitiveData(fields map[string]any) map[string]any {
	result := make(map[string]any)
	for key, value := range fields {
		switch key {
		case "password", "token", "secret", "access_token", "refresh_token":
			result[key] = "***"
		default:
			result[key] = value
		}
	}
	return result
}

type errorWithLogCtx struct {
	next error
	ctx  logCtx
}

func (e *errorWithLogCtx) Error() string {
	return e.next.Error()
}

func WrapError(ctx context.Context, err error) error {
	var e *errorWithLogCtx
	if errors.As(err, &e) {
		return err
	}
	c := logCtx{}
	if x, ok := ctx.Value(slogFields).(logCtx); ok {
		c = x
	}
	return &errorWithLogCtx{
		next: err,
		ctx:  c,
	}
}

func OriginalError(err error) error {
	var e *errorWithLogCtx
	if errors.As(err, &e) {
		return e.next
	}
	return err
}

func ErrorCtx(ctx context.Context, err error) context.Context {
	var e *errorWithLogCtx
	if errors.As(err, &e) {
		return context.WithValue(ctx, slogFields, e.ctx)
	}
	return ctx
}
