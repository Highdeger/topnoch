package xlog

func LogTrace(msg string) {
	logThis(logTrace, msg)
}

func LogDebug(msg string) {
	logThis(logDebug, msg)
}

func LogInfo(msg string) {
	logThis(logInfo, msg)
}

func LogWarning(msg string) {
	logThis(logWarning, msg)
}

func LogError(msg string) {
	logThis(logError, msg)
}

func LogFatal(msg string) {
	logThis(logFatal, msg)
}

func LogPanic(msg string) {
	logThis(logPanic, msg)
}
