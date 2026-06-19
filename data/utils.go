package data

import (
	"fmt"
	"time"
)

func RosLogMessage(topicName string, level LogLevel, rawMessage interface{}) LogMessage {
    return LogMessage{
        Timestamp: time.Now(),
        Level: level,
        Source: fmt.Sprintf("ROS:%s", topicName),
        Payload: fmt.Sprintf("%v", rawMessage),
    }
}

func YoloLogMessage(cameraID string, level LogLevel, rawMessage interface{}) LogMessage {
    return LogMessage{
        Timestamp: time.Now(),
        Level: level,
        Source: fmt.Sprintf("YOLO:%s", cameraID),
        Payload: fmt.Sprintf("%v", rawMessage),
    }
}

func DockLogMessage(containerID string, level LogLevel, rawMessage interface{}) LogMessage {
    return LogMessage{
        Timestamp: time.Now(),
        Level: level,
        Source: fmt.Sprintf("DOCK:%s", containerID),
        Payload: fmt.Sprintf("%v", rawMessage),
    }
}
