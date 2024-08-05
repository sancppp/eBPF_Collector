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

package userspace

import (
	"log"

	"ebpf_exporter/event"

	"ebpf_exporter/userspace/cnetwork"
	"ebpf_exporter/userspace/csyscall"
	"ebpf_exporter/userspace/fileopen"

	"github.com/cilium/ebpf/rlimit"
)

// ebpf的userspace部分

var (
	stopChannels = make([]chan struct{}, 0)
)

func registerEbpf(init func(<-chan struct{}, chan<- event.IEvent), eventCh chan<- event.IEvent) {
	stopper := make(chan struct{}, 1)
	go init(stopper, eventCh)
	stopChannels = append(stopChannels, stopper)
}

func Run(stopper <-chan struct{}, eventCh chan<- event.IEvent) {
	// Allow the current process to lock memory for eBPF resources.
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal(err)
	}

	//注册ebpf程序

	registerEbpf(fileopen.InitFileopen, eventCh)
	registerEbpf(cnetwork.InitCNetwork, eventCh)
	registerEbpf(csyscall.InitSyscall, eventCh)

	<-stopper
	for _, stopCh := range stopChannels {
		stopCh <- struct{}{}
	}
}
