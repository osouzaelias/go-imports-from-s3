package telemetry

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var instance trace.Tracer

func GetTracer() trace.Tracer {
	if instance == nil {
		instance = otel.Tracer("go-import-from-s3")
	}
	return instance
}
