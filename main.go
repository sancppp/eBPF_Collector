package main

import (
	"ebpf_exporter/comsumer"
	"ebpf_exporter/event"
	"ebpf_exporter/userspace"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	//全局ebpf事件管道
	eventCh = make(chan event.IEvent, 1000)
)

func main() {
	stopper := make(chan os.Signal, 1)
	signal.Notify(stopper, os.Interrupt, syscall.SIGTERM)

	ebpfstopper := make(chan struct{}, 1)

	go userspace.Run(ebpfstopper, eventCh)
	// go route.Init()
	// go comsumer.StartPrometheus(eventCh)
	// go comsumer.StartLog(eventCh)
	// go comsumer.StartPrint(eventCh)
	go comsumer.StartHttp(eventCh)
	<-stopper
	ebpfstopper <- struct{}{}
	time.Sleep(1 * time.Second)
}
