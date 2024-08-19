package logger

import (
	"context"
	"time"

	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
	"google.golang.org/grpc"
)

func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	logger.Custom.Infoln("Received grpc request: ", info.FullMethod)
	startTime := time.Now()
	status := "ok"

	resp, err := handler(ctx, req)
	if err != nil {
		status = "ko"
	}

	logger.Custom.Infoln(
		"uri", info.FullMethod,
		"method", "grpc",
		"duration", time.Since(startTime),
		"status", status,
	)

	return resp, err
}
