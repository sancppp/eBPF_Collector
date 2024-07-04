package event

type UdpSendmsg_event struct {
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

func (UdpSendmsg_event) GetName() string {
	return "UdpSendmsg_event"
}

func (e UdpSendmsg_event) GetTimestamp() int64 {
	return e.Timestamp
}

func (e UdpSendmsg_event) GetPid() uint32 {
	return e.Pid
}

func (e UdpSendmsg_event) GetComm() string {
	return e.Comm
}
