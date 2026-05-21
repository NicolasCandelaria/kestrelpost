package ending

import "strconv"

// Ending IDs match docs/superpowers/specs/2026-05-14-kestrel-post-ending-evaluator-design.md
type Ending uint8

const (
	TheRelay       Ending = 1
	DarkFrequency  Ending = 2
	TheKidWasRight Ending = 3
	FullBroadcast  Ending = 4
	TheConvoy      Ending = 5
	DeadAir        Ending = 6
	Fallback       Ending = 7
)

func (e Ending) String() string {
	switch e {
	case TheRelay:
		return "THE_RELAY"
	case DarkFrequency:
		return "DARK_FREQUENCY"
	case TheKidWasRight:
		return "THE_KID_WAS_RIGHT"
	case FullBroadcast:
		return "FULL_BROADCAST"
	case TheConvoy:
		return "THE_CONVOY"
	case DeadAir:
		return "GONE_DARK"
	case Fallback:
		return "FALLBACK"
	default:
		return "Ending(" + strconv.Itoa(int(e)) + ")"
	}
}
