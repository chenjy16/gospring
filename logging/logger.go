package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// Logger defines the interface used for logging GoSpring events.
// Implementations of this interface can customize how GoSpring logs its internal events.
type Logger interface {
	// LogEvent logs the given event
	LogEvent(event Event)
}

// NopLogger is a Logger that ignores all events.
// Use this when you want to disable GoSpring's internal logging.
var NopLogger Logger = &nopLogger{}

type nopLogger struct{}

func (l *nopLogger) LogEvent(event Event) {
	// Do nothing
}

// ConsoleLogger is a Logger that writes human-readable messages to the console.
// This is the default logger used by GoSpring.
type ConsoleLogger struct {
	// W is the writer to write logs to. Defaults to os.Stderr.
	W io.Writer
	
	// mu protects concurrent writes to W
	mu sync.Mutex
}

// NewConsoleLogger creates a new ConsoleLogger that writes to os.Stderr.
func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{
		W: os.Stderr,
	}
}

// NewConsoleLoggerWithWriter creates a new ConsoleLogger that writes to the specified writer.
func NewConsoleLoggerWithWriter(w io.Writer) *ConsoleLogger {
	return &ConsoleLogger{
		W: w,
	}
}

// LogEvent logs the given event to the console in a human-readable format.
func (l *ConsoleLogger) LogEvent(event Event) {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	if l.W == nil {
		l.W = os.Stderr
	}
	
	fmt.Fprintln(l.W, event.String())
}

// StandardLogger wraps Go's standard log.Logger to implement the Logger interface.
type StandardLogger struct {
	logger *log.Logger
}

// NewStandardLogger creates a new StandardLogger using the provided log.Logger.
func NewStandardLogger(logger *log.Logger) *StandardLogger {
	return &StandardLogger{
		logger: logger,
	}
}

// NewStandardLoggerWithPrefix creates a new StandardLogger with the specified prefix.
func NewStandardLoggerWithPrefix(prefix string) *StandardLogger {
	return &StandardLogger{
		logger: log.New(os.Stderr, prefix, log.LstdFlags),
	}
}

// LogEvent logs the given event using the standard logger.
func (l *StandardLogger) LogEvent(event Event) {
	l.logger.Println(event.String())
}

// MultiLogger allows logging to multiple loggers simultaneously.
type MultiLogger struct {
	loggers []Logger
}

// NewMultiLogger creates a new MultiLogger that logs to all provided loggers.
func NewMultiLogger(loggers ...Logger) *MultiLogger {
	return &MultiLogger{
		loggers: loggers,
	}
}

// LogEvent logs the given event to all configured loggers.
func (l *MultiLogger) LogEvent(event Event) {
	for _, logger := range l.loggers {
		logger.LogEvent(event)
	}
}

// AddLogger adds a logger to the MultiLogger.
func (l *MultiLogger) AddLogger(logger Logger) {
	l.loggers = append(l.loggers, logger)
}

// FilteredLogger wraps another logger and only logs events that match the filter function.
type FilteredLogger struct {
	logger Logger
	filter func(Event) bool
}

// NewFilteredLogger creates a new FilteredLogger that only logs events matching the filter.
func NewFilteredLogger(logger Logger, filter func(Event) bool) *FilteredLogger {
	return &FilteredLogger{
		logger: logger,
		filter: filter,
	}
}

// LogEvent logs the event only if it passes the filter.
func (l *FilteredLogger) LogEvent(event Event) {
	if l.filter(event) {
		l.logger.LogEvent(event)
	}
}

// LeveledLogger provides different log levels for different event types.
type LeveledLogger struct {
	logger Logger
	level  LogLevel
}

// LogLevel represents the logging level.
type LogLevel int

const (
	// LogLevelDebug logs all events including detailed debugging information
	LogLevelDebug LogLevel = iota
	// LogLevelInfo logs informational events (default level)
	LogLevelInfo
	// LogLevelWarn logs warning and error events only
	LogLevelWarn
	// LogLevelError logs error events only
	LogLevelError
	// LogLevelOff disables all logging
	LogLevelOff
)

// NewLeveledLogger creates a new LeveledLogger with the specified level.
func NewLeveledLogger(logger Logger, level LogLevel) *LeveledLogger {
	return &LeveledLogger{
		logger: logger,
		level:  level,
	}
}

// LogEvent logs the event based on the configured log level.
func (l *LeveledLogger) LogEvent(event Event) {
	eventLevel := l.getEventLevel(event)
	if eventLevel >= l.level {
		l.logger.LogEvent(event)
	}
}

// getEventLevel determines the log level for a given event.
func (l *LeveledLogger) getEventLevel(event Event) LogLevel {
	switch event.(type) {
	case *DependencyInjectionFailed:
		return LogLevelError
	case *LifecycleStarted:
		if e := event.(*LifecycleStarted); e.Error != nil {
			return LogLevelError
		}
		return LogLevelInfo
	case *LifecycleStopped:
		if e := event.(*LifecycleStopped); e.Error != nil {
			return LogLevelError
		}
		return LogLevelInfo
	case *ComponentScanned, *DependencyInjected:
		return LogLevelDebug
	case *ComponentRegistered, *ComponentCreated, *ComponentDestroyed:
		return LogLevelInfo
	case *LifecycleStarting, *LifecycleStopping:
		return LogLevelDebug
	case *ContextStarting, *ContextStarted, *ContextStopping, *ContextStopped:
		return LogLevelInfo
	case *ContainerCreated:
		return LogLevelInfo
	default:
		return LogLevelInfo
	}
}