package main

import (
	"context"
	"log"
	"os"

	geometry_msgs "github.com/Voltamon/ros-router/msgs/geometry_msgs/msg"
	std_msgs "github.com/Voltamon/ros-router/msgs/std_msgs/msg"

	"github.com/tiiuae/rclgo/pkg/rclgo"
	"github.com/tiiuae/rclgo/pkg/rclgo/types"
)

func Ros() {
  args, _, err := rclgo.ParseArgs(os.Args[1:])
  if err != nil {
    log.Fatalf("Failed to parse args: %s", err.Error())
  }

  if err := rclgo.Init(args); err != nil {
    log.Fatalf("Failed to initialize ROS 2: %s", err.Error())
  }
  defer rclgo.Uninit()

  node, err := rclgo.NewNode("middleware_logger", "default")
  if err != nil {
    log.Fatalf("Failed to create node: %s", err.Error())
  }

  topicsToTrack := []struct {
    Name string
    TypeSupport types.MessageTypeSupport
  }{
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
      }
      log.Printf("%s", msg)
    })
  }

  if err != nil {
    log.Fatalf("Failed to create subscription: %s", err.Error())
  }

  if err = rclgo.Spin(context.Background()); err != nil {
    log.Fatalf("Failed to spin: %s", err.Error())
  }
}
