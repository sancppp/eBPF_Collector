package comsumer

import (
	"ebpf_exporter/event"
	"encoding/json"
	"log"
	"os"
)

func StartLog(eventCh <-chan event.IEvent) {
	//重定向日志到文件
	f, err := os.OpenFile("ebpf_exporter.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("open log file error: %v", err)
	}
	defer f.Close()
	//清空文件
	f.Truncate(0)
	log.SetOutput(f)

	for event_ := range eventCh {
		// 将类型+json作为日志打印
		eventJson, err := json.Marshal(event_)
		if err != nil {
			log.Printf("json.Marshal error: %v", err)
			continue
		}
		log.Printf("event type: %T, event: %s", event_, eventJson)
	}
}

func StartPrint(eventCh <-chan event.IEvent) {
	for event_ := range eventCh {
		// 将类型+json作为日志打印
		eventJson, err := json.Marshal(event_)
		if err != nil {
			log.Printf("json.Marshal error: %v", err)
			continue
		}
		log.Printf("event type: %T, event: %s", event_, eventJson)
	}
}
