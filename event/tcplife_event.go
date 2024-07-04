package event

type Tcplife_event struct {
	Timestamp int64
	Saddr     [16]byte /* uint128 */
	Daddr     [16]byte /* uint128 */
	TsUs      uint64
	SpanUs    uint64
	RxB       uint64
	TxB       uint64
	Pid       uint32
	Sport     uint16
	Dport     uint16
	Family    uint16
	Comm      string
}

func (Tcplife_event) GetName() string {
	return "tcplife_event"
}

func (e Tcplife_event) GetTimestamp() int64 {
	return e.Timestamp
}

func (e Tcplife_event) GetPid() uint32 {
	return e.Pid
}

func (e Tcplife_event) GetComm() string {
	return e.Comm
}
