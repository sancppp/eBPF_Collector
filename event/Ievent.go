package event

type IEvent interface {
	GetTimestamp() int64
	GetName() string
	GetPid() uint32
	GetComm() string
}
