package userspace

import (
	"log"

	"ebpf_exporter/event"

	"github.com/cilium/ebpf/rlimit"
)

// ebpf的userspace部分

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go bpf ../kern/tcplife.c -- -I../kern/headers

var (
	EventCh = make(chan event.IEvent, 100)
)

func Run(stopper <-chan struct{}) {
	// Allow the current process to lock memory for eBPF resources.
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal(err)
	}
	tcpliftstoper := make(chan struct{}, 1)
	initTcpLife(tcpliftstoper)

	<-stopper
	tcpliftstoper <- struct{}{}
}
