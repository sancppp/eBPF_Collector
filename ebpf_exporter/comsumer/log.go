package comsumer

import (
	"ebpf_exporter/event"
	"encoding/json"
	"log"
	"net/http"
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
func StartHttp(eventCh <-chan event.IEvent) {

	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		events := make([]event.IEvent, 0, 200)
		for i := 0; i < 200; i++ {
			select {
			case event := <-eventCh:
				events = append(events, event)
			default:
				// 如果通道中没有事件，立即返回
				goto END
			}
		}
	END:
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(events); err != nil {
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		}
		defer r.Body.Close()
	})
	if err := http.ListenAndServe("192.168.0.202:8089", nil); err != nil {
		log.Println("http server err: ", err)
	}

}
