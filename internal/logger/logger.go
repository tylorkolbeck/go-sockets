package logger

import "log"

type Logger struct {
	prefix string
}

func NewLogger(prefix string) *Logger {
	return &Logger{
		prefix: prefix,
	}
}

func (l *Logger) Info(msg string, args ...any) {
	log.Printf("[%s] INFO: "+msg, append([]any{l.prefix}, args...)...)
}

func (l *Logger) Error(msg string, args ...any) {
	log.Printf("[%s] ERROR: "+msg, append([]any{l.prefix}, args...))
}
