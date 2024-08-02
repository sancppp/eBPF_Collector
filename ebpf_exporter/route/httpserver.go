package route

import (
	"ebpf_exporter/util"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func containerInfoHandler(w http.ResponseWriter, r *http.Request) {
	containerInfos, err := util.GetContainerInfo()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting container info: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(containerInfos)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func ContainerInfoServer() {
	http.HandleFunc("/containerinfo", containerInfoHandler)
	log.Println("Starting containerinfo server on http://192.168.0.202:8888/containerinfo")
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
