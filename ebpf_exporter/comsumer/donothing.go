package comsumer

import "ebpf_exporter/event"

func DoNothing(eventCh <-chan event.IEvent) {
	for range eventCh {
		// do nothing
	}
}
