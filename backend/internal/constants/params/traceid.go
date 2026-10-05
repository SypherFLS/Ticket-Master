package params

import (
	"context"
	"log/slog"
	"tmaster/internal/constants"
)

func GetRequestID(ctx context.Context) string {
	requestID, _ := ctx.Value(constants.RequestIDKey).(string)

	return requestID
}

func GetLogger(ctx context.Context) *slog.Logger {
	logger, _ := ctx.Value(constants.RequestLogger).(*slog.Logger)

	return logger
}
