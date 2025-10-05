package repo

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type CEPRepository struct {
	weatherApiHost string
	weatherApiPort string
}

func NewCEPRepository(WEATHER_API_host string, WEATHER_API_port string) *CEPRepository {
	return &CEPRepository{
		weatherApiHost: WEATHER_API_host,
		weatherApiPort: WEATHER_API_port,
	}
}

func (r *CEPRepository) IsValid(cep_address string) bool {
	check, _ := regexp.MatchString("^[0-9]{8}$", cep_address)
	return (len(cep_address) == 8 && cep_address != "" && check)
}

func (r *CEPRepository) Get(cep_address string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf(
		"http://%s:%s/cep/%s",
		r.weatherApiHost,
		r.weatherApiPort,
		cep_address),
		nil,
	)
	if err != nil {
		log.Printf("Fail to create the request: %v", err)
		return err
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport,
			otelhttp.WithSpanNameFormatter(func(_ string, req *http.Request) string {
				return "get-cep-temp"
			}),
		),
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Fail to make the request: %v", err)
		return err
	}
	defer resp.Body.Close()

	ctx_err := ctx.Err()
	if ctx_err != nil {
		select {
		case <-ctx.Done():
			err := ctx.Err()
			log.Printf("Max timeout reached: %v", err)
			return err
		}
	}

	return nil
}
