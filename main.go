package main

import (
	"ebpf_exporter/event"
	"ebpf_exporter/prometheus"
	"ebpf_exporter/route"
	"ebpf_exporter/userspace"
	"os"
	"os/signal"
	"syscall"
)

var (
	//全局ebpf事件管道
	eventCh = make(chan event.IEvent, 1000)
)

func main() {
	stopper := make(chan os.Signal, 1)
	signal.Notify(stopper, os.Interrupt, syscall.SIGTERM)

	ebpfstopper := make(chan struct{}, 1)
	defer func() { ebpfstopper <- struct{}{} }()

	go userspace.Run(ebpfstopper, eventCh)
	go route.Init()
	go prometheus.Comsumer(eventCh)

	<-stopper
}
