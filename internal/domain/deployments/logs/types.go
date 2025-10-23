package logs

import "time"

type LogLevel string

const (
	LogLevelInfo    LogLevel = "info"
	LogLevelWarn    LogLevel = "warn"
	LogLevelError   LogLevel = "error"
	LogLevelSuccess LogLevel = "success"
)

type StructuredLog struct {
	Timestamp       string   `json:"timestamp"`
	RelativeSeconds int      `json:"relative_seconds"`
	Level           LogLevel `json:"level"`
	Message         string   `json:"message"`
	Raw             string   `json:"raw"`
}

type LogParser struct {
	startTime time.Time
}

func NewLogParser(startTime time.Time) *LogParser {
	return &LogParser{
		startTime: startTime,
	}
}
