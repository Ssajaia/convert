package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Level represents a log severity level.
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger writes structured log lines to one or more writers.
type Logger struct {
	mu      sync.Mutex
	writers []io.Writer
	level   Level
}

var std = &Logger{
	writers: []io.Writer{os.Stderr},
	level:   INFO,
}

// Init configures the global logger. Call once at startup.
func Init(logFile string, level Level) error {
	std.mu.Lock()
	defer std.mu.Unlock()

	std.level = level
	std.writers = []io.Writer{os.Stderr}

	if logFile != "" {
		f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("open log file: %w", err)
		}
		std.writers = append(std.writers, f)
	}
	return nil
}

func log(level Level, format string, args ...any) {
	std.mu.Lock()
	defer std.mu.Unlock()

	if level < std.level {
		return
	}

	msg := fmt.Sprintf(format, args...)
	line := fmt.Sprintf("%s [%s] %s\n", time.Now().Format("2006-01-02T15:04:05"), level, msg)

	for _, w := range std.writers {
		_, _ = io.WriteString(w, line)
	}
}

func Debug(format string, args ...any) { log(DEBUG, format, args...) }
func Info(format string, args ...any)  { log(INFO, format, args...) }
func Warn(format string, args ...any)  { log(WARN, format, args...) }
func Error(format string, args ...any) { log(ERROR, format, args...) }
