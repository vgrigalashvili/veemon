package mqtt

import (
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var client mqtt.Client

// Connect connects to the MQTT broker using the given brokerURL and clientID.
func Connect(brokerURL, clientID string) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerURL)
	opts.SetClientID(clientID)
	opts.SetCleanSession(true)
	opts.SetConnectTimeout(30 * time.Second)

	opts.OnConnect = func(c mqtt.Client) {
		log.Println("[MQTT] Connected to broker")
	}
	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Printf("[MQTT] Connection lost: %v", err)
	}

	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error connecting to MQTT broker: %v", token.Error())
	}
	return client
}

// SubscribeHeartbeat subscribes to the heartbeat topic.
func SubscribeHeartbeat(topic string) {
	if token := client.Subscribe(topic, 1, heartbeatHandler); token.Wait() && token.Error() != nil {
		log.Fatalf("Error subscribing to heartbeat topic: %v", token.Error())
	}
	log.Printf("Subscribed to heartbeat topic: %s", topic)
}

func heartbeatHandler(client mqtt.Client, msg mqtt.Message) {
	// Handle the heartbeat message (for example, log it).
	log.Printf("Heartbeat received on topic %s: %s", msg.Topic(), msg.Payload())
}
