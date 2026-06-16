package main

import (
	"log"
	"os"

	std_msgs "github.com/Voltamon/ros-router/msgs/std_msgs/msg"
	"github.com/tiiuae/rclgo/pkg/rclgo"
)

func ros() {
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

  _, err = node.NewSubscription("/topic",std_msgs.StringTypeSupport, nil, func (sub *rclgo.Subscription) {
    var msg std_msgs.String
    _, err = sub.TakeMessage(&msg)
    if err != nil {
      log.Printf("Failed to take message: %s", err.Error())
      return
    }
  })
}
