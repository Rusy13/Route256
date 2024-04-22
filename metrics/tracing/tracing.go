package postgresql

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"io"
	"log"
)

func initTracer(service string) (opentracing.Tracer, io.Closer, error) {
	// Инициализация конфигурации Jaeger
	cfg, err := jaegercfg.FromEnv()
	if err != nil {
		log.Printf("Could not parse Jaeger env vars: %s", err.Error())
		return nil, nil, err
	}

	// Установка имени сервиса
	cfg.ServiceName = service

	// Инициализация трейсера Jaeger
	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		log.Printf("Could not initialize jaeger tracer: %s", err.Error())
		return nil, nil, err
	}

	// Установка глобального трейсера
	opentracing.SetGlobalTracer(tracer)

	return tracer, closer, nil
}
