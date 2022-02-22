package SQL

import "fmt"

type LogType int64

const (
	LogInfo    LogType = iota
	LogWarning
	LogError
	LogFatal
)

func (e LogType) String() string {
	switch e {
	case LogInfo:
		return "INFO"
	case LogWarning:
		return "WARNING"
	case LogError:
		return "ERROR"
	case LogFatal:
		return "FATAL"
	default:
		return fmt.Sprintf("%d", e)
	}
}
