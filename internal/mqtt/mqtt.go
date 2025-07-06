package mqtt

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	client mqtt.Client
}

func NewMQTTClient(broker string, clientID string) *MQTTClient {
	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID(clientID).
		SetCleanSession(true).
		SetAutoReconnect(true)

	client := mqtt.NewClient(opts)
	return &MQTTClient{client: client}
}

func (m *MQTTClient) Connect() error {
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
	}
	return nil
}

func (m *MQTTClient) Subscribe(topic string, qos byte, handler func([]byte)) error {
	messageHandler := func(client mqtt.Client, msg mqtt.Message) {
		handler(msg.Payload())
	}

	if token := m.client.Subscribe(topic, qos, messageHandler); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %w", topic, token.Error())
	}
	return nil
}

func (m *MQTTClient) Disconnect() {
	m.client.Disconnect(250)
}
