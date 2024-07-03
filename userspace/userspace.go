package userspace

import (
	"log"

	"ebpf_exporter/event"

	tcplife "ebpf_exporter/userspace/tcp_life"
	udprcv "ebpf_exporter/userspace/udp_rcv"

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
	registerEbpf(tcplife.InitTcpLife, eventCh)
	registerEbpf(udprcv.InitUdpRcv, eventCh)

	<-stopper
	for _, stopCh := range stopChannels {
		stopCh <- struct{}{}
	}
}
