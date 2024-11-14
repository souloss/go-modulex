package telemetry

import (
	"context"

	"github.com/souloss/go-clean-arch/pkg/version"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// InitOpenTelemetry 初始化 OpenTelemetry 全局 Tracer
func InitOpenTelemetry() error {
	// 配置 OTLP 导出器
	// exp, err := otlptrace.New(
	// 	context.Background(),
	// 	otlptracegrpc.NewClient(),
	// )
	exp, err := stdouttrace.New()
	if err != nil {
		return err
	}

	// 配置 Resource
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", version.ProjectName),
		),
	)
	if err != nil {
		return err
	}

	// 配置 Tracer Provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return nil
}
