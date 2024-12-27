package shared

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type InterceptorFunc func(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	hander grpc.UnaryHandler,
) (
	interface{},
	error,
)

var vlog slog.Logger

func CreateLogInterceptor(log slog.Logger) InterceptorFunc {
	vlog = log
	return logger
}

func logger(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	hander grpc.UnaryHandler,
) (
	interface{},
	error,
) {
	uuid, _ := uuid.NewV7()
	vlog := vlog.With("id", uuid.String())

	vlog.Debug("msg", info.FullMethod, req)
	res, err := hander(ctx, req)
	if err != nil {
		vlog.Error(err.Error())
	} else {
		vlog.Debug("ok", "resp", res)
	}

	return res, err
}
