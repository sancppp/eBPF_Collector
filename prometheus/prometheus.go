package prometheus

import (
	"ebpf_exporter/event"
	"fmt"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// 过期时间
const TIMEOUT_S = 60

var (
	// 指标
	tcpByteCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ebpf_tcp_byte_counter",
			Help: "进程视角一次tcp发包、收包事件",
		},
		[]string{"timestamp", "protocol", "flag", "pid", "comm", "daddr", "dport", "saddr", "sport"},
	)
	tcpByteSum = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ebpf_tcp_byte_sum",
			Help: "进程视角tcp发包、收包的字节数和",
		},
		[]string{"protocol", "flag", "pid", "comm"},
	)
	udpByteCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ebpf_udp_byte_counter",
			Help: "进程视角一次udp发包、收包事件",
		},
		[]string{"timestamp", "protocol", "flag", "pid", "comm", "daddr", "dport", "saddr", "sport"},
	)
	udpByteSum = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ebpf_udp_byte_sum",
			Help: "进程视角udp发包、收包的字节数和",
		},
		[]string{"protocol", "flag", "pid", "comm"},
	)
)

func init() {
	// 注册指标
	prometheus.MustRegister(tcpByteCounter)
	prometheus.MustRegister(tcpByteSum)
	prometheus.MustRegister(udpByteCounter)
	prometheus.MustRegister(udpByteSum)
}

func Comsumer(eventCh <-chan event.IEvent) {
	for event_ := range eventCh {
		// 将数据作为指标发送到Prometheus exporter
		switch event_ := event_.(type) {
		case event.Udp_event:
			{
				tmp := prometheus.Labels{
					"timestamp": fmt.Sprintf("%d", event_.GetTimestamp()),
					"protocol":  "udp",
					"flag":      fmt.Sprintf("%d", event_.Flag),
					"pid":       fmt.Sprintf("%d", event_.GetPid()),
					"comm":      event_.GetComm(),
					"daddr":     fmt.Sprintf("%v", event_.Daddr),
					"dport":     fmt.Sprintf("%d", event_.Dport),
					"saddr":     fmt.Sprintf("%v", event_.Saddr),
					"sport":     fmt.Sprintf("%d", event_.Sport),
				}
				udpByteCounter.With(tmp).Add(float64(event_.Len))
				go func(tmp *prometheus.Labels) {
					time.Sleep(TIMEOUT_S * time.Second)
					udpByteCounter.Delete(*tmp)
				}(&tmp)
				udpByteSum.With(prometheus.Labels{
					"protocol": "udp",
					"flag":     fmt.Sprintf("%d", event_.Flag),
					"pid":      fmt.Sprintf("%d", event_.GetPid()),
					"comm":     event_.GetComm(),
				}).Add(float64(event_.Len))
			}
		case event.Tcp_event:
			{
				tmp := prometheus.Labels{
					"timestamp": fmt.Sprintf("%d", event_.GetTimestamp()),
					"protocol":  "tcp",
					"flag":      fmt.Sprintf("%d", event_.Flag),
					"pid":       fmt.Sprintf("%d", event_.GetPid()),
					"comm":      event_.GetComm(),
					"daddr":     fmt.Sprintf("%v", event_.Daddr),
					"dport":     fmt.Sprintf("%d", event_.Dport),
					"saddr":     fmt.Sprintf("%v", event_.Saddr),
					"sport":     fmt.Sprintf("%d", event_.Sport),
				}
				tcpByteCounter.With(tmp).Add(float64(event_.Len))
				go func(tmp *prometheus.Labels) {
					time.Sleep(TIMEOUT_S * time.Second)
					tcpByteCounter.Delete(*tmp)
				}(&tmp)
				tcpByteSum.With(prometheus.Labels{
					"protocol": "tcp",
					"flag":     fmt.Sprintf("%d", event_.Flag),
					"pid":      fmt.Sprintf("%d", event_.GetPid()),
					"comm":     event_.GetComm(),
				}).Add(float64(event_.Len))

			}
		default:
			{
				log.Printf("unknown event type: %T", event_)
			}
		}
	}
}
