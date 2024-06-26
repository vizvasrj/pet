package middleware

import (
	"context"
	"runtime/debug"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RecoveryMiddleware(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = status.Errorf(codes.Internal, "internal server error: %v", r)
				logger.Error("Recovered from panic",
					zap.Any("panic", r),
					zap.String("method", info.FullMethod),
					zap.String("stack", string(debug.Stack())),
				)
			}
		}()
		return handler(ctx, req)
	}
}
