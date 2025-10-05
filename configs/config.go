package configs

import (
	"log"

	"github.com/spf13/viper"
)

type Conf struct {
	cepApiHttpPort                        string `mapstructure:"CEP_API_HTTP_PORT"`
	cepApiOtelServiceName                 string `mapstructure:"CEP_API_OTEL_SERVICE_NAME"`
	weatherApiPort                        string `mapstructure:"WEATHER_API_PORT"`
	weatherApiHost                        string `mapstructure:"WEATHER_API_HOST"`
	OpenWeathermapApiKey                  string `mapstructure:"OPEN_WEATHERMAP_API_KEY"`
	weatherApiServiceName                 string `mapstructure:"WEATHER_API_SERVICE_NAME"`
	OpenTelemetryCollectorExporerEndpoint string `mapstructure:"OPEN_TELEMETRY_COLLECTOR_EXPORTER_ENDPOINT"`
}

func LoadConfig(path string) (*Conf, error) {
	var cfg *Conf
	var err error

	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()
	viper.BindEnv("CEP_API_HTTP_PORT")
	viper.BindEnv("CEP_API_OTEL_SERVICE_NAME")
	viper.BindEnv("WEATHER_API_PORT")
	viper.BindEnv("WEATHER_API_HOST")
	viper.BindEnv("OPEN_WEATHERMAP_API_KEY")
	viper.BindEnv("WEATHER_API_SERVICE_NAME")
	viper.BindEnv("OPEN_TELEMETRY_COLLECTOR_EXPORTER_ENDPOINT")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("WARNING: %v\n", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return cfg, err
}
