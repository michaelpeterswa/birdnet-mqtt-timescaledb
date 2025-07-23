package mqtt

import (
	"fmt"
	"log/slog"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type SubscriptionInfo struct {
	qos     byte
	handler func([]byte)
}

type MQTTClient struct {
	client        mqtt.Client
	subscriptions map[string]SubscriptionInfo
}

type MQTTClientOption func(*mqtt.ClientOptions)

func WithMQTTClientID(clientID string) MQTTClientOption {
	return func(opts *mqtt.ClientOptions) {
		opts.SetClientID(clientID)
	}
}

func WithMQTTKeepAlive(keepAlive time.Duration) MQTTClientOption {
	return func(opts *mqtt.ClientOptions) {
		opts.SetKeepAlive(keepAlive)
	}
}

func WithMQTTPingTimeout(pingTimeout time.Duration) MQTTClientOption {
	return func(opts *mqtt.ClientOptions) {
		opts.SetPingTimeout(pingTimeout)
	}
}

func WithMQTTConnectTimeout(connectTimeout time.Duration) MQTTClientOption {
	return func(opts *mqtt.ClientOptions) {
		opts.SetConnectTimeout(connectTimeout)
	}
}

func WithMQTTCleanSession(cleanSession bool) MQTTClientOption {
	return func(opts *mqtt.ClientOptions) {
		opts.SetCleanSession(cleanSession)
	}
}

func NewMQTTClient(broker string, options ...MQTTClientOption) *MQTTClient {
	mqttClient := &MQTTClient{
		subscriptions: make(map[string]SubscriptionInfo),
	}

	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetAutoReconnect(true).
		SetConnectionLostHandler(func(client mqtt.Client, err error) {
			slog.Error("MQTT connection lost", slog.String("error", err.Error()))
		}).
		SetOnConnectHandler(func(client mqtt.Client) {
			slog.Info("MQTT connection established")
			mqttClient.resubscribeAll(client)
		}).
		SetReconnectingHandler(func(client mqtt.Client, opts *mqtt.ClientOptions) {
			slog.Warn("MQTT attempting to reconnect")
		})

	for _, option := range options {
		option(opts)
	}

	mqttClient.client = mqtt.NewClient(opts)
	return mqttClient
}

func (m *MQTTClient) Connect() error {
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
	}
	return nil
}

func (m *MQTTClient) resubscribeAll(client mqtt.Client) {
	for topic, info := range m.subscriptions {
		messageHandler := func(client mqtt.Client, msg mqtt.Message) {
			info.handler(msg.Payload())
		}
		
		if token := client.Subscribe(topic, info.qos, messageHandler); token.Wait() && token.Error() != nil {
			slog.Error("failed to re-subscribe to topic", 
				slog.String("topic", topic), 
				slog.String("error", token.Error().Error()))
		} else {
			slog.Info("re-subscribed to topic", slog.String("topic", topic))
		}
	}
}

func (m *MQTTClient) Subscribe(topic string, qos byte, handler func([]byte)) error {
	messageHandler := func(client mqtt.Client, msg mqtt.Message) {
		handler(msg.Payload())
	}

	if token := m.client.Subscribe(topic, qos, messageHandler); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %w", topic, token.Error())
	}

	m.subscriptions[topic] = SubscriptionInfo{
		qos:     qos,
		handler: handler,
	}
	
	slog.Info("subscribed to topic", slog.String("topic", topic))
	return nil
}

func (m *MQTTClient) Disconnect() {
	m.client.Disconnect(250)
}
