package route

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Init() {
	http.Handle("/metrics", promhttp.Handler())
	log.Println("http server start")

	if err := http.ListenAndServe(":2112", nil); err != nil {
		log.Println("route err: ", err)
	}
}
