package fileopen

import (
	"bytes"
	"ebpf_exporter/event"
	"ebpf_exporter/util"
	"encoding/binary"
	"errors"
	"log"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
	"golang.org/x/sys/unix"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -target amd64 -type fileopen_event bpf ../../kern/fileopen.c -- -I../../kern/headers

func InitFileopen(stopper <-chan struct{}, eventCh chan<- event.IEvent) {
	// Load pre-compiled programs and maps into the kernel.
	objs := bpfObjects{}
	if err := loadBpfObjects(&objs, nil); err != nil {
		log.Fatalf("loading objects: %v", err)
	}
	defer objs.Close()

	kp, err := link.Kretprobe("do_filp_open", objs.KretprobeDoFilpOpen, nil)
	if err != nil {
		log.Fatalf("attaching raw tracepoint: %v", err)
	}
	defer kp.Close()

	rd, err := ringbuf.NewReader(objs.FileopenRb)
	if err != nil {
		log.Fatalf("opening ringbuf reader: %s", err)
	}
	defer rd.Close()
	log.Println("Waiting for fileopen events..")

	go func() {
		// bpfEvent is generated by bpf2go.
		var bpfevent bpfFileopenEvent
		for {
			record, err := rd.Read()
			if err != nil {
				if errors.Is(err, ringbuf.ErrClosed) {
					log.Println("Received signal, exiting..")
					return
				}
				log.Printf("reading from reader: %s", err)
				continue
			}
			// Parse the ringbuf event entry into a bpfEvent structure.
			if err := binary.Read(bytes.NewBuffer(record.RawSample), binary.LittleEndian, &bpfevent); err != nil {
				log.Printf("parsing ringbuf event: %s", err)
				continue
			}
			filename := ""
			for i := 32 - 1; i >= 0; i-- {
				if bpfevent.Filename[i][0] != 0 {
					filename += unix.ByteSliceToString(bpfevent.Filename[i][:])
				}
			}
			// 打印 bpfEvent 结构体的内容
			fileopenEvent := event.Fileopen_event{
				Type:          "Fileopen_event",
				Timestamp:     bpfevent.Timestamp,
				Pid:           bpfevent.Pid,
				Comm:          unix.ByteSliceToString(bpfevent.Comm[:]),
				Filename:      filename,
				Fsname:        unix.ByteSliceToString(bpfevent.Fsname[:]),
				Cid:           unix.ByteSliceToString(bpfevent.Cid[:]),
				ContainerName: util.GetContainerName(unix.ByteSliceToString(bpfevent.Cid[:])),
			}
			eventCh <- fileopenEvent
		}
	}()

	<-stopper
}
