package logger

import (
	"io"
	"log"
	"os"
	"sync"
)

type LogLevel int

const (
	LevelError LogLevel = iota
	LevelInfo
	LevelDebug
)

var (
	DebugLogger *log.Logger
	ErrorLogger *log.Logger
	InfoLogger  *log.Logger

	currentLevel LogLevel
	levelMutex   sync.RWMutex
)

func init() {
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLogger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	SetLogLevel(LevelInfo) // Default to INFO level
}

func SetLogLevel(level LogLevel) {
	levelMutex.Lock()
	defer levelMutex.Unlock()
	currentLevel = level

	// ErrorLogger always remains active
	
	if level < LevelInfo {
		InfoLogger.SetOutput(ioutil.Discard)
	} else {
		InfoLogger.SetOutput(os.Stdout)
	}

	if level < LevelDebug {
		DebugLogger.SetOutput(ioutil.Discard)
	} else {
		DebugLogger.SetOutput(os.Stdout)
	}
}

func GetLogLevel() LogLevel {
	levelMutex.RLock()
	defer levelMutex.RUnlock()
	return currentLevel
}

func IsDebugEnabled() bool {
    return GetLogLevel() >= LevelDebug
}

func IsInfoEnabled() bool {
    return GetLogLevel() >= LevelInfo
}