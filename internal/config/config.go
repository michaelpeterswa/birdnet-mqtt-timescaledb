package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"error"`

	MQTTAddress        string        `env:"MQTT_ADDRESS" envDefault:"tcp://localhost:1883"`
	MQTTClientID       string        `env:"MQTT_CLIENT_ID" envDefault:"birdnet-mqtt-timescaledb"`
	MQTTTopic          string        `env:"MQTT_TOPIC" envDefault:"birdnet"`
	MQTTQoS            byte          `env:"MQTT_QOS" envDefault:"1"`
	MQTTKeepAlive      time.Duration `env:"MQTT_KEEPALIVE" envDefault:"30s"`
	MQTTPingTimeout    time.Duration `env:"MQTT_PING_TIMEOUT" envDefault:"10s"`
	MQTTConnectTimeout time.Duration `env:"MQTT_CONNECT_TIMEOUT" envDefault:"30s"`
	MQTTCleanSession   bool          `env:"MQTT_CLEAN_SESSION" envDefault:"true"`

	TimescaleConnString string `env:"TIMESCALE_CONN_STRING" envDefault:"postgres://postgres:example@timescaledb:5432/postgres?sslmode=disable"`

	MetricsEnabled bool `env:"METRICS_ENABLED" envDefault:"true"`
	MetricsPort    int  `env:"METRICS_PORT" envDefault:"8081"`

	Local bool `env:"LOCAL" envDefault:"false"`

	Timezone string `env:"TIMEZONE" envDefault:"America/Los_Angeles"`

	TracingEnabled    bool    `env:"TRACING_ENABLED" envDefault:"false"`
	TracingSampleRate float64 `env:"TRACING_SAMPLERATE" envDefault:"0.01"`
	TracingService    string  `env:"TRACING_SERVICE" envDefault:"katalog-agent"`
	TracingVersion    string  `env:"TRACING_VERSION"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	err := env.Parse(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}
