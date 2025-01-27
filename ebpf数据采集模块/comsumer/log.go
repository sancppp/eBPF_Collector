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
