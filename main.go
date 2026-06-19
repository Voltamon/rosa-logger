package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Voltamon/ros-router/data"
	"github.com/Voltamon/ros-router/service"
	"github.com/tiiuae/rclgo/pkg/rclgo"

	geometry_msgs "github.com/Voltamon/ros-router/msgs/geometry_msgs/msg"
	std_msgs "github.com/Voltamon/ros-router/msgs/std_msgs/msg"
	worker "github.com/Voltamon/ros-router/service/workers"
)

func main() {
    logStream := make(chan data.LogMessage, 100)
    ctx, cancel := context.WithCancel(context.Background())
    waitGroup := &sync.WaitGroup{}

    args, _, err := rclgo.ParseArgs(os.Args[1:])
    if err != nil {
        log.Fatalf("Failed to parse ROS 2 arguments: %s", err.Error())
    }

    waitGroup.Add(1)
    go func() {
        defer waitGroup.Done()

        topicsToTrack := []data.TopicConfig{
{"/chatter", std_msgs.StringTypeSupport},
{"/cmd_vel", geometry_msgs.TwistTypeSupport},
        }

        err := worker.StartRosWorker(ctx, args, topicsToTrack, logStream)
        if err != nil {
            log.Fatalf("ROS Worker error: %s", err.Error())
        }
    }()

    // waitGroup.Add(1)
    // go func() {
    //     defer waitGroup.Done()

    //     err := worker.StartDockWorker(ctx, []string{"blissful_brattain"}, logStream)
    //     if err != nil {
    //         log.Fatalf("Docker Worker error: %s", err.Error())
    //     }
    // }()

    waitGroup.Add(1)
    go func() {
        defer waitGroup.Done()

        err := service.StartGRPCServer(":50051", logStream)
        if err != nil {
            log.Fatalf("gRPC Server error: %s", err.Error())
        }
    }()

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan

    cancel()
    waitGroup.Wait()
    close(logStream)
}
