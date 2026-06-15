package main

import (
	"log"
	"os"

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

  _, err = rclgo.NewNode("middleware_logger", "default")
  if err != nil {
    log.Fatalf("Failed to create node: %s", err.Error())
  }
}
