package logger

import (
	"io"
	"log"
	"os"
)

func New(file string, debug bool) *Loggers {
	fileHandle := getFileHandle(file)
	loggers := Loggers{
		Error: newLogger("ERROR", fileHandle, os.Stderr),
	}
	if debug {
		loggers.Debug = newLogger("DEBUG", fileHandle, os.Stdout)
	} else {
		loggers.Debug = log.New(io.Discard, "", 0)
	}
	return &loggers
}

func getFileHandle(file string) *os.File {
	if file != "" {
		fileHandle, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}
		return fileHandle
	}
	return nil
}

func newLogger(prefix string, file *os.File, alternative *os.File) *log.Logger {
	if file != nil {
		return log.New(file, prefix+": ", log.Ldate|log.Ltime)
	}
	return log.New(alternative, prefix+": ", log.Ldate|log.Ltime)
}

type Loggers struct {
	Debug *log.Logger
	Error *log.Logger
}
