package event

type Tcp_event struct {
	Timestamp int64
	Flag      uint16 // 0:send,1:recv
	Pid       uint32
	Daddr     [16]byte /* uint128 */
	Dport     uint16
	Saddr     [16]byte /* uint128 */
	Sport     uint16
	Len       uint16
	Comm      string
}

func (Tcp_event) GetName() string {
	return "Tcp_event"
}

func (e Tcp_event) GetTimestamp() int64 {
	return e.Timestamp
}

func (e Tcp_event) GetPid() uint32 {
	return e.Pid
}

func (e Tcp_event) GetComm() string {
	return e.Comm
}
