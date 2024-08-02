package event

type IEvent interface {
	GetTimestamp() uint64
	GetName() string
	GetPid() uint32
	GetComm() string
}
