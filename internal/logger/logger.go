package logger

import (
	"context"
	"log/slog"
)

var logger *slog.Logger = slog.Default()

func Init(l *slog.Logger) {
	logger = l
}

func With(args ...any) *slog.Logger {
	return logger.With(args...)
}

func WithGroup(name string) *slog.Logger {
	return logger.WithGroup(name)
}

func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

func Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	logger.Log(ctx, level, msg, args...)
}

func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}

func Enabled(ctx context.Context, level slog.Level) bool {
	return logger.Enabled(ctx, level)
}

func LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	logger.LogAttrs(ctx, level, msg, attrs...)
}

func DebugContext(ctx context.Context, msg string, args ...any) {
	logger.DebugContext(ctx, msg, args...)
}

func InfoContext(ctx context.Context, msg string, args ...any) {
	logger.InfoContext(ctx, msg, args...)
}

func WarnContext(ctx context.Context, msg string, args ...any) {
	logger.WarnContext(ctx, msg, args...)
}

func ErrorContext(ctx context.Context, msg string, args ...any) {
	logger.ErrorContext(ctx, msg, args...)
}
