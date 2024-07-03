package event

type Tcplife_event struct {
	Saddr  [16]byte /* uint128 */
	Daddr  [16]byte /* uint128 */
	TsUs   uint64
	SpanUs uint64
	RxB    uint64
	TxB    uint64
	Pid    uint32
	Sport  uint16
	Dport  uint16
	Family uint16
	Comm   string
}

func (Tcplife_event) Name() string {
	return "tcplife_event"
}
