package event

type CNetwork_event struct {
	Type          string
	Timestamp     uint64
	Pid           uint32
	Comm          string
	Cid           string
	ContainerName string
	Flag          uint8
	Daddr         [4]byte
	Dport         uint16
	Saddr         [4]byte
	Sport         uint16
}

func (CNetwork_event) GetName() string {
	return "CNetwork_event"
}

func (e CNetwork_event) GetTimestamp() uint64 {
	return e.Timestamp
}

func (e CNetwork_event) GetPid() uint32 {
	return e.Pid
}

func (e CNetwork_event) GetComm() string {
	return e.Comm
}
