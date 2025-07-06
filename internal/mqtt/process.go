package mqtt

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/michaelpeterswa/birdnet-mqtt-timescaledb/internal/birdnet"
	"github.com/michaelpeterswa/birdnet-mqtt-timescaledb/internal/timescale"
)

type Processor struct {
	mqttClient      *MQTTClient
	timescaleClient *timescale.TimescaleClient
	ctx             context.Context
}

func NewProcessor(mqttClient *MQTTClient, timescaleClient *timescale.TimescaleClient, ctx context.Context) *Processor {
	return &Processor{
		mqttClient:      mqttClient,
		timescaleClient: timescaleClient,
		ctx:             ctx,
	}
}

func (p *Processor) Start(topic string, qos byte) error {
	return p.mqttClient.Subscribe(topic, qos, p.processMessage)
}

func (p *Processor) processMessage(payload []byte) {
	var detection birdnet.BirdDetection
	if err := json.Unmarshal(payload, &detection); err != nil {
		slog.Error("failed to unmarshal bird detection", slog.String("error", err.Error()))
		return
	}

	event, err := detection.ToBirdDetectionEvent()
	if err != nil {
		slog.Error("failed to convert bird detection to event", slog.String("error", err.Error()))
		return
	}

	if err := p.timescaleClient.StoreBirdDetectionEvent(p.ctx, event); err != nil {
		slog.Error("failed to store bird detection event", slog.String("error", err.Error()))
		return
	}

	slog.Info("stored bird detection event",
		slog.String("common_name", event.CommonName),
		slog.String("scientific_name", event.ScientificName),
		slog.Float64("confidence", event.Confidence))
}
