package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/MadAppGang/dingo/pkg/dgo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/trace"
)

type Config struct {
	Endpoint   string `json:"endpoint"`
	Insecure   bool   `json:"insecure"`
	TestMarker string `json:"test_marker"`
}

func loadConfig(filepath string) dgo.Result[*Config, error] {
	tmp, err := os.ReadFile(filepath)
	if err != nil {
		return dgo.Err[*Config](err)
	}
	data := tmp

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return dgo.Err[*Config](err)
	}

	return dgo.Ok[*Config, error](&config)
}

func initTracer(ctx context.Context, endpoint string, insecure bool) dgo.Result[*trace.TracerProvider, error] {
	tmp1, err1 := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err1 != nil {
		return dgo.Err[*trace.TracerProvider](err1)
	}
	exporter := tmp1

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)

	fmt.Printf("✓ OpenTelemetry tracer initialized with OTLP GRPC exporter. Endpoint: %s\n", endpoint)

	return dgo.Ok[*trace.TracerProvider, error](tp)
}

func testTraces(ctx context.Context, testMarker string) error {
	tracer := otel.Tracer("test-tracer")
	_, span := tracer.Start(ctx, "test-span")
	defer span.End()

	span.SetAttributes(
		attribute.String("gen_ai.completion.0.content", "Fake response prompt"),
		attribute.String("gen_ai.prompt.0.content", "Fake request prompt"),
		attribute.String("test.marker", testMarker),
	)

	fmt.Println("✓ OpenTelemetry traces test completed")
	fmt.Println("If you go to the 'Distributed Tracing' application in Dynatrace and switch from 'requests' to 'spans' on the top left hand side of the screen, you should see the trace there.")
	return nil
}

func main() {
	configPath := flag.String("config", "config.json", "Path to configuration file")
	flag.Parse()

	var config *Config
	res1 := loadConfig(*configPath)
	if res1.IsOk() {
		c := res1.MustOk()
		{
			config = c
		}
	} else {
		err := res1.MustErr()
		{
			log.Fatalf("Failed to load config: %v", err)
		}
	}

	ctx := context.Background()

	var tp *trace.TracerProvider
	res := initTracer(ctx, config.Endpoint, config.Insecure)
	if res.IsOk() {
		provider := res.MustOk()
		{
			tp = provider
		}
	} else {
		err := res.MustErr()
		{
			log.Fatalf("Failed to initialize tracer: %v", err)
		}
	}

	defer tp.Shutdown(ctx)

	if err := testTraces(ctx, config.TestMarker); err != nil {
		log.Fatalf("Traces test failed: %v", err)
	}
}
