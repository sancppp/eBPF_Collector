// Copyright (C) 2024 Tianzhenxiong
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package comsumer

import (
	"context"
	"encoding/json"
	"log"

	"ebpf_exporter/event"

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
		js, _ := json.Marshal(evt)
		_, err := conn.WriteMessages(
			kafka.Message{Value: js},
		)
		if err != nil {
			log.Printf("could not write message: %v", err)
		}
	}
}
