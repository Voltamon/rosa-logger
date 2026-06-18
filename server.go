package main

import (
	"log"
	"net"
	"time"

	"github.com/Voltamon/ros-router/buffer"
	"google.golang.org/grpc"
)

type LogServer struct {
    buffer.UnimplementedLogServiceServer
    logStream <-chan LogMessage
}

func (server* LogServer) StreamLogs(request *buffer.LogStreamRequest, stream buffer.LogService_StreamLogsServer) error {
    log.Println("[gRPC] New client connected")

    for msg := range server.logStream {
        pubMsg := &buffer.LogStream{
            Timestamp: msg.Timestamp.Format(time.RFC3339),
            Level: string(msg.Level),
            Source: msg.Source,
            Payload: msg.Payload,
        }

        if err := stream.Send(pubMsg); err != nil {
            log.Printf("[gRPC] Client disconnected: %s", err.Error())
        }
    }

  return nil
}

func StartGRPCServer(port string, globalLogStream <-chan LogMessage) error {
    listener, err := net.Listen("tcp", port)
    if err != nil {
        return err
    }

    grpcServer := grpc.NewServer()
    buffer.RegisterLogServiceServer(grpcServer, &LogServer{
        logStream: globalLogStream,
    })

    log.Printf("[gRPC] Server listening on %s", port)
    return grpcServer.Serve(listener)
}
