package log

import (
	"os"
)

type Log struct {
	debugLogFile *os.File
}

func CreateLog(file string) *Log {
	debugLogFile, err := os.Create(file)

	if err != nil {
		panic("")
	}

	return &Log{
		debugLogFile: debugLogFile,
	}
}

func (l *Log) Debug(msg string) {
	l.debugLogFile.WriteString(msg + "\n")
}
