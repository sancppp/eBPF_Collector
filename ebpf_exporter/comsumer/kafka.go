package comsumer

import (
	"context"
	"log"

	"ebpf_exporter/event"

	"github.com/bytedance/sonic"
	"github.com/segmentio/kafka-go"
)

var (
	addr      = "192.168.0.249:19092"
	topic     = "event"
	partition = 0
)

func StartKafka(eventCh <-chan event.IEvent) {

	conn, err := kafka.DialLeader(context.Background(), "tcp", addr, topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}
	// Close the writer when done
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatal("failed to close writer:", err)
		}
	}()

	// Start a goroutine to read from eventCh and write to Kafka
	for evt := range eventCh {
		js, _ := sonic.Marshal(evt)
		_, err := conn.WriteMessages(
			kafka.Message{Value: js},
		)
		if err != nil {
			log.Printf("could not write message: %v", err)
		}
	}
}
