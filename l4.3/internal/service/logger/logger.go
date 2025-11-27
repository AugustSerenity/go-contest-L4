package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/AugustSerenity/go-contest-L4/l4.3_Events-calendar/internal/model"
)

type Logger struct {
	logChan chan model.LogEntry
	done    chan bool
}

func New() *Logger {
	return &Logger{
		logChan: make(chan model.LogEntry, 100),
		done:    make(chan bool),
	}
}

func (l *Logger) Start() {
	go func() {
		for entry := range l.logChan {
			jsonEntry, _ := json.Marshal(entry)
			fmt.Fprintf(os.Stdout, "%s\n", jsonEntry)
		}
		l.done <- true
	}()
}

func (l *Logger) Stop() {
	close(l.logChan)
	<-l.done
}

func (l *Logger) Log(level, message, requestID string) {
	entry := model.LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		RequestID: requestID,
	}

	select {
	case l.logChan <- entry:
	default:
		jsonEntry, _ := json.Marshal(entry)
		fmt.Fprintf(os.Stdout, "FALLBACK: %s\n", jsonEntry)
	}
}

func (l *Logger) Info(message, requestID string) {
	l.Log("INFO", message, requestID)
}

func (l *Logger) Error(message, requestID string) {
	l.Log("ERROR", message, requestID)
}

func (l *Logger) Warn(message, requestID string) {
	l.Log("WARN", message, requestID)
}
