package event

type UdpSendmsg_event struct {
	Timestamp int64
	Flag      uint16
	Pid       uint32
	Daddr     [16]byte /* uint128 */
	Dport     uint16
	Saddr     [16]byte /* uint128 */
	Sport     uint16
	Len       uint16
	Comm      string
}

func (UdpSendmsg_event) Name() string {
	return "UdpSendmsg_event"
}
