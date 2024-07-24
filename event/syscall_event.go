package event

type Syscall_event struct {
	Timestamp uint64
	Flag      uint8 // 0 for enter, 1 for exit, 2 for count
	Pid       uint32
	Comm      string
	Syscall   uint32
	Ret       int64
	Cid       string
	Info      string
}

func (Syscall_event) GetName() string {
	return "System_event"
}

func (e Syscall_event) GetTimestamp() uint64 {
	return e.Timestamp
}

func (e Syscall_event) GetPid() uint32 {
	return e.Pid
}

func (e Syscall_event) GetComm() string {
	return e.Comm
}
