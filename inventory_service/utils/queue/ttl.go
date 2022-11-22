package queue

type TtlLevel uint16

const (
	Level1 TtlLevel = iota // 1s
	Level2                 //  5s
	Level3                 // 10s
	Level4                 // 30s
	Level5                 // 1min
	Level6                 // 5min
	Level7                 // 30 min
)

func (l TtlLevel) String() string {
	switch l {
	case Level1: // 1s
		return "1s"
	case Level2: //  5s
		return "2s"
	case Level3: // 10s
		return "10s"
	case Level4: // 30s
		return "30s"
	case Level5: // 1min
		return "1min"
	case Level6: // 5min
		return "5min"
	case Level7: // 30 min
		return "30min"
	}
	return ""
}
