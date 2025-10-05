package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/JoaoPedroVicentin/observabilidade-open-telemetry/cep-api/internal/infra/web"
	"github.com/JoaoPedroVicentin/observabilidade-open-telemetry/cep-api/internal/infra/web/webserver"
	"github.com/JoaoPedroVicentin/observabilidade-open-telemetry/configs"
	otel_provider "github.com/JoaoPedroVicentin/observabilidade-open-telemetry/pkg/otel"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
)

func ConfigureServer(conf *configs.Conf) *webserver.WebServer {
	fmt.Println("Starting web server on port", conf.cepApiHttpPort)

	tracer := otel.Tracer("intput-api-tracer")

	webserver := webserver.NewWebServer(":" + conf.cepApiHttpPort)
	webCEPHandler := web.NewWebCEPHandler(conf, tracer)
	webStatusHandler := web.NewWebStatusHandler()
	webserver.AddHandler("POST /cep", webCEPHandler.Get)
	webserver.AddHandler("GET /status", webStatusHandler.Get)
	return webserver
}

func init() {
	viper.AutomaticEnv()
}

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := otel_provider.InitProvider(configs.cepApiOtelServiceName, configs.OpenTelemetryCollectorExporerEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	go func() {
		webserver := ConfigureServer(configs)
		webserver.Start()
	}()

	select {
	case <-sigCh:
		log.Println("Shutting down gracefully, CTRL+c pressed...")
	case <-ctx.Done():
		log.Println("Shutting down due other reason...")
	}

	_, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
}
