package graylog

const (
	LevelEmergency     Level = 0
	LevelAlert         Level = 1
	LevelCritical      Level = 2
	LevelError         Level = 3
	LevelWarning       Level = 4
	LevelNotice        Level = 5
	LevelInformational Level = 6
	LevelDebug         Level = 7
	LevelUnkown        Level = 8
)

type Level int

func (lv Level) String() string {
	switch lv {
	case LevelEmergency:
		return "EMERG"
	case LevelAlert:
		return "ALERT"
	case LevelCritical:
		return "CRITI"
	case LevelError:
		return "ERROR"
	case LevelWarning:
		return "WARNI"
	case LevelNotice:
		return "NOTIC"
	case LevelInformational:
		return "INFOR"
	case LevelDebug:
		return "DEBUG"
	}

	return "-----"
}
