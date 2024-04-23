package interceptor

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
)

// TracingUnaryServerInterceptor создает трейсинг для unary запросов gRPC.
func TracingUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Создаем новый спан для запроса
		span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, opentracing.GlobalTracer(), info.FullMethod)
		defer span.Finish()

		// Выполняем gRPC-метод
		resp, err := handler(ctx, req)

		// Добавляем метаданные к спану
		if err != nil {
			ext.LogError(span, err)
		}
		span.LogKV("method", info.FullMethod, "error", err)

		return resp, err
	}
}
