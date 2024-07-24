package event

type Fileopen_event struct {
	Timestamp uint64
	Pid       uint32
	Comm      string
	Filename  string
	Fsname    string
	Cid       string
}

func (Fileopen_event) GetName() string {
	return "Fileopen_event"
}

func (e Fileopen_event) GetTimestamp() uint64 {
	return e.Timestamp
}

func (e Fileopen_event) GetPid() uint32 {
	return e.Pid
}

func (e Fileopen_event) GetComm() string {
	return e.Comm
}
