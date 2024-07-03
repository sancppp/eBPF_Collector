package route

import (
	"ebpf_exporter/event"
	"ebpf_exporter/userspace"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Init() {
	startPrometheushttpserver()
	startNetEventHttpServer()
	fmt.Println("http server start")
	http.ListenAndServe(":2112", nil)
}

func startPrometheushttpserver() {
	http.Handle("/metrics", promhttp.Handler())
}

func startNetEventHttpServer() {
	http.Handle("/tcplife", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 返回当前eventchan中的数据,json格式
		eventList := make([]event.IEvent, 0)
		for {
			select {
			case event := <-userspace.EventCh:
				eventList = append(eventList, event)
			default:
				goto END
			}
		}
	END:
		data, err := json.Marshal(eventList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}))
}
