package prometheus

import (
	"ebpf_exporter/event"
	"encoding/json"
	"fmt"
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// 创建一个带标签的 CounterVec
	tcpRxEventCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ebpf_tcp_rx_event",
			Help: "tcp rx event counter , rx bytes 目前的逻辑只会记录一次完成的TCP活动",
		},
		[]string{"timestamp", "protocol", "pid", "comm", "daddr", "dport", "saddr", "sport"},
	)
	tcpTxEventCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ebpf_tcp_tx_event",
			Help: "tcp tx event counter , tx bytes 目前的逻辑只会记录一次完成的TCP活动",
		},
		[]string{"timestamp", "protocol", "pid", "comm", "daddr", "dport", "saddr", "sport"},
	)

	tcpRxByteSum = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ebpf_tcp_rx_byte_sum",
			Help: "进程视角的TCP接收字节数",
		},
		[]string{"protocol", "pid", "comm"},
	)

	tcpTxByteSum = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ebpf_tcp_tx_byte_sum",
			Help: "进程视角的TCP发送字节数",
		},
		[]string{"protocol", "pid", "comm"},
	)

	udpEventCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ebpf_udp_event",
			Help: "udp event counter , bytes",
		},
		[]string{"timestamp", "protocol", "flag", "pid", "comm", "daddr", "dport", "saddr", "sport"},
	)
)

func init() {
	// 注册指标
	prometheus.MustRegister(tcpRxEventCounter)
	prometheus.MustRegister(tcpTxEventCounter)
	prometheus.MustRegister(tcpRxByteSum)
	prometheus.MustRegister(tcpTxByteSum)
	prometheus.MustRegister(udpEventCounter)
}

func Comsumer(eventCh <-chan event.IEvent) {
	for event_ := range eventCh {
		jsonstr, _ := json.Marshal(event_)
		log.Print(string(jsonstr))
		// 将数据作为指标发送到Prometheus exporter
		switch event_ := event_.(type) {
		case event.Tcplife_event:
			{
				tcpRxEventCounter.With(prometheus.Labels{
					"timestamp": fmt.Sprintf("%d", event_.GetTimestamp()),
					"protocol":  "tcp",
					"pid":       fmt.Sprintf("%d", event_.GetPid()),
					"comm":      event_.GetComm(),
					"daddr":     fmt.Sprintf("%v", event_.Daddr),
					"dport":     fmt.Sprintf("%d", event_.Dport),
					"saddr":     fmt.Sprintf("%v", event_.Saddr),
					"sport":     fmt.Sprintf("%d", event_.Sport),
				}).Add(float64(event_.RxB))
				tcpTxEventCounter.With(prometheus.Labels{
					"timestamp": fmt.Sprintf("%d", event_.GetTimestamp()),
					"protocol":  "tcp",
					"pid":       fmt.Sprintf("%d", event_.GetPid()),
					"comm":      event_.GetComm(),
					"daddr":     fmt.Sprintf("%v", event_.Daddr),
					"dport":     fmt.Sprintf("%d", event_.Dport),
					"saddr":     fmt.Sprintf("%v", event_.Saddr),
					"sport":     fmt.Sprintf("%d", event_.Sport),
				}).Add(float64(event_.TxB))
				tcpRxByteSum.With(prometheus.Labels{
					"protocol": "tcp",
					"pid":      fmt.Sprintf("%d", event_.GetPid()),
					"comm":     event_.GetComm(),
				}).Add(float64(event_.RxB))
				tcpTxByteSum.With(prometheus.Labels{
					"protocol": "tcp",
					"pid":      fmt.Sprintf("%d", event_.GetPid()),
					"comm":     event_.GetComm(),
				}).Add(float64(event_.TxB))
			}
		case event.UdpSendmsg_event:
			{
				udpEventCounter.With(prometheus.Labels{
					"timestamp": fmt.Sprintf("%d", event_.GetTimestamp()),
					"protocol":  "udp",
					"flag":      fmt.Sprintf("%d", event_.Flag),
					"pid":       fmt.Sprintf("%d", event_.GetPid()),
					"comm":      event_.GetComm(),
					"daddr":     fmt.Sprintf("%v", event_.Daddr),
					"dport":     fmt.Sprintf("%d", event_.Dport),
					"saddr":     fmt.Sprintf("%v", event_.Saddr),
					"sport":     fmt.Sprintf("%d", event_.Sport),
				}).Add(float64(event_.Len))
			}
		default:
			{
				log.Printf("unknown event type: %T", event_)
			}
		}
	}
}
