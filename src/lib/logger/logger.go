package logger

import "log"

type Logger interface {
	Debug(message string)
	Info(message string)
	Error(message string)
}

type StdoutLogger struct{}

func (l *StdoutLogger) Debug(message string) {
	log.Println("DEBUG: " + message)
}

func (l *StdoutLogger) Info(message string) {
	log.Println("INFO:  " + message)
}

func (l *StdoutLogger) Error(message string) {
	log.Println("ERROR: " + message)
}
