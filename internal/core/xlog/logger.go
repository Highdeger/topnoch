package xlog

import (
	"../xfile"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type logType int

const (
	logTrace logType = iota
	logDebug
	logInfo
	logWarning
	logError
	logFatal
	logPanic
)

var (
	LogPrefixes = []string{"TRACE: ", "DEBUG: ", "INFO: ", "WARNING: ", "ERROR: ", "FATAL: ", "PANIC: "}
)

func logDir() string {
	dir := "coreLog"
	//TODO: update coreLog directory path
	return dir
}

func logFlags() int {
	flags := log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix
	//TODO: update flags from database
	return flags
}

func logThis(typ logType, msg string) {
	t := time.Now()
	name := fmt.Sprintf("%4d-m%02d-d%02d.coreLog", t.Year(), t.Month(), t.Day())
	if !xfile.DirExists(logDir()) {
		if err := os.Mkdir(logDir(), os.ModePerm); err != nil {
			panic(err)
		}
	}
	absFilePath, err := filepath.Abs(logDir() + string(os.PathSeparator) + name)
	if err != nil {
		panic(err)
	}
	file, e := os.OpenFile(absFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if e != nil {
		panic(e)
	}
	logger := log.New(file, LogPrefixes[typ], logFlags())
	logger.Println(msg)
	logger = nil
	if err = file.Close(); err != nil {
		panic(err)
	}
}
