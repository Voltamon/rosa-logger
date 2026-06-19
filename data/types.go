package data

import (
	"time"

	"github.com/tiiuae/rclgo/pkg/rclgo/types"
)
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

type TopicConfig struct {
    Name string
    TypeSupport types.MessageTypeSupport
}
