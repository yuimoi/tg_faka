package my_log

import (
	"log"
)

type LogLevel int

const (
	LogLevelDebugValue = 0
	LogLevelInfoValue  = 1
	LogLevelWarnValue  = 2
	LogLevelErrorValue = 3
)

func Log(level int, message string) {
	if level >= getLogLevel() {
		log.Println(message)
	}
}

func LogDebug(message string) {
	Log(LogLevelDebugValue, message)
}
func LogInfo(message string) {
	Log(LogLevelInfoValue, message)
}
func LogWarn(message string) {
	Log(LogLevelWarnValue, message)
}
func LogError(message string) {
	Log(LogLevelErrorValue, message)
}
func getLogLevel() int {

	return 1
}
