package logger

import (
	"log"
	"os"
)

// Logger handles application logging
type Logger struct {
	infoLog  *log.Logger
	warnLog  *log.Logger
	errorLog *log.Logger
}

// NewLogger creates a new logger instance
func NewLogger() *Logger {
	return &Logger{
		infoLog:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		warnLog:  log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLog: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Info logs an info message
func (l *Logger) Info(message string, args ...interface{}) {
	if len(args) > 0 {
		l.infoLog.Printf(message, args...)
	} else {
		l.infoLog.Println(message)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(message string, err error) {
	if err != nil {
		l.warnLog.Printf("%s: %v", message, err)
	} else {
		l.warnLog.Println(message)
	}
}

// Error logs an error message
func (l *Logger) Error(message string, err error) {
	if err != nil {
		l.errorLog.Printf("%s: %v", message, err)
	} else {
		l.errorLog.Println(message)
	}
}

// Fatal logs a fatal error and exits
func (l *Logger) Fatal(message string, err error) {
	if err != nil {
		l.errorLog.Fatalf("%s: %v", message, err)
	} else {
		l.errorLog.Fatalln(message)
	}
}
