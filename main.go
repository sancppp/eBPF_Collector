package main

import (
	"ebpf_exporter/route"
	"ebpf_exporter/userspace"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	stopper := make(chan os.Signal, 1)
	signal.Notify(stopper, os.Interrupt, syscall.SIGTERM)
	ebpfstopper := make(chan struct{}, 1)
	go userspace.Run(ebpfstopper)
	go route.Init()

	<-stopper
	ebpfstopper <- struct{}{}
}
