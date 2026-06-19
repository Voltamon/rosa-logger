package worker

import (
	"context"
	"fmt"
	"log"

	"github.com/Voltamon/ros-router/data"
	geometry_msgs "github.com/Voltamon/ros-router/msgs/geometry_msgs/msg"
	std_msgs "github.com/Voltamon/ros-router/msgs/std_msgs/msg"

	"github.com/tiiuae/rclgo/pkg/rclgo"
)

func StartRosWorker(ctx context.Context, args *rclgo.Args, logChan chan <- data.LogMessage) error {
    if err := rclgo.Init(args); err != nil {
        fmt.Errorf("Failed to initialize ROS 2: %s", err.Error())
    }
    defer rclgo.Uninit()

    node, err := rclgo.NewNode("middleware_logger", "default")
    if err != nil {
        return fmt.Errorf("Failed to create node: %s", err.Error())
    }

    topicsToTrack := []data.TopicConfig{
    {"/chatter", std_msgs.StringTypeSupport},
    {"/cmd_vel", geometry_msgs.TwistTypeSupport},
    }

    for _, config := range topicsToTrack {
        topicName := config.Name
        typeSupport := config.TypeSupport

        _, err = node.NewSubscription(topicName, typeSupport, nil, func (sub *rclgo.Subscription) {
            msg := typeSupport.New()
            _, err := sub.TakeMessage(msg)

            if err != nil {
                log.Printf("Failed to take message from topic %s: %s", topicName, err.Error())
                return
            }
            logChan <- data.RosLogMessage(topicName, data.LevelInfo, msg)
        })
    }

    if err != nil {
        return fmt.Errorf("Failed to create subscription: %s", err.Error())
    }

    return rclgo.Spin(ctx)
}
