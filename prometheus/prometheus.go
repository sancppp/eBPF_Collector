package prometheus

import (
	"ebpf_exporter/event"
	"encoding/json"
	"log"
)

func Comsumer(eventCh <-chan event.IEvent) {
	for event := range eventCh {
		jsonstr, _ := json.Marshal(event)
		log.Print(string(jsonstr))
	}
}
