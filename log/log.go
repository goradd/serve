package log

import (
	"context"
	"log/slog"
)

const LogLevelServeDebug = -5

func Debug(ctx context.Context, msg string, args ...any) {
	slog.Default().Log(ctx, LogLevelServeDebug, msg, args...)
}
