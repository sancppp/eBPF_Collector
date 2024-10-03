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
	log.Println("Starting containerinfo server on http://192.168.252.131:8888/containerinfo")
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
