package SQL

type LogType int64

const(
	LogInfo LogType = 1
	LogWarning = 2
	LogError = 3
	LogFatal = 4
)
