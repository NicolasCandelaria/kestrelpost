package night

type Mode byte

const (
	ModePreShift Mode = iota
	ModeReceive
	ModeScan
	ModeIncident
	ModeLogbook
)

func (m Mode) String() string {
	switch m {
	case ModePreShift:
		return "PRE_SHIFT"
	case ModeReceive:
		return "RECEIVE"
	case ModeScan:
		return "SCAN"
	case ModeIncident:
		return "INCIDENT"
	case ModeLogbook:
		return "LOGBOOK"
	default:
		return "UNKNOWN"
	}
}

func NextAfterChoice(_ int) Mode {
	return ModeIncident
}
