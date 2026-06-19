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
	worker "github.com/Voltamon/ros-router/service/workers"
	"github.com/tiiuae/rclgo/pkg/rclgo"
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

        err := worker.StartRosWorker(ctx, args, logStream)
        if err != nil {
            log.Fatalf("ROS Worker error: %s", err.Error())
        }
    }()

    waitGroup.Add(1)
    go func() {
        defer waitGroup.Done()

        err := worker.StartDockWorker(ctx, "blissful_brattain", logStream)
        if err != nil {
            log.Fatalf("Docker Worker error: %s", err.Error())
        }
    }()

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
