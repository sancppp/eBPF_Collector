package cmd

import (
	"ebpf_exporter/comsumer"
	"ebpf_exporter/event"
	"ebpf_exporter/route"
	"ebpf_exporter/userspace"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var eventCh = make(chan event.IEvent, 1000)

var rootCmd = &cobra.Command{
	Use:   "ebpf_exporter",
	Short: "ebpf采集器",
	Long:  `西电mobisys实验室容器云团队，基于ebpf技术的容器行为监控系统采集器`,
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	stopper := make(chan os.Signal, 1)
	signal.Notify(stopper, os.Interrupt, syscall.SIGTERM)

	ebpfstopper := make(chan struct{}, 1)

	go userspace.Run(ebpfstopper, eventCh)

	// 启动http服务, 容器信息
	go route.ContainerInfoServer()

	// 启动消费者
	// go comsumer.StartPrint(eventCh)
	// go comsumer.StartHttp(eventCh)
	// go comsumer.DoNothing(eventCh)
	go comsumer.StartKafka(eventCh)
	<-stopper
	ebpfstopper <- struct{}{}
	time.Sleep(1 * time.Second)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
