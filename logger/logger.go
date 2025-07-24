package logger

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogKey string

const LogContext LogKey = "LogContext"

func NewZap(level string) (*zap.Logger, func()) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.LevelKey = "level"
	encoderConfig.NameKey = "logger"
	encoderConfig.CallerKey = "caller"
	encoderConfig.MessageKey = "message"
	encoderConfig.StacktraceKey = "stacktrace"

	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(parseLevel(level)),
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := cfg.Build()
	if err != nil {
		panic("failed to build logger: " + err.Error())
	}

	return logger, func() {
		_ = logger.Sync()
	}
}

func parseLevel(level string) zapcore.Level {
	level = strings.ToLower(level)
	switch level {
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	case "debug":
		return zapcore.DebugLevel
	default:
		return zapcore.InfoLevel
	}
}

func FromContext(ctx context.Context) (*zap.Logger, error) {
	l, ok := ctx.Value(LogContext).(*zap.Logger)
	if !ok {
		return nil, errors.New("unable get log from context")
	}
	return l, nil
}

func ZapLoggerMiddleware(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			stop := time.Now()

			fields := []zap.Field{
				zap.String("method", c.Request().Method),
				zap.String("path", c.Request().URL.Path),
				zap.Int("status", c.Response().Status),
				zap.Duration("latency", stop.Sub(start)),
				zap.String("remote_ip", c.RealIP()),
				zap.String("user_agent", c.Request().UserAgent()),
			}

			if err != nil {
				logger.Error("http request error", append(fields, zap.Error(err))...)
			} else {
				logger.Info("http request", fields...)
			}

			return err
		}
	}
}

func LoggerToContextMiddleware(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := context.WithValue(c.Request().Context(), LogContext, logger)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
