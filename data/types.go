package data

import "time"
type LogLevel string

const (
    LevelInfo LogLevel = "INFO"
    LevelWarn LogLevel = "WARN"
    LevelError LogLevel = "ERROR"
)

type LogMessage struct {
    Timestamp time.Time
    Level     LogLevel
    Source    string
    Payload   string
}
