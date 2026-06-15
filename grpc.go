package main

import (
	"context"
	"log"
	"net"

	buf "github.com/Voltamon/ros-router/buf"
	"google.golang.org/grpc"
)

type server struct {
  buf.UnimplementedPingServiceServer
}

func (s *server) SendPing(ctx context.Context, req *buf.PingRequest) (*buf.PingResponse, error) {
  log.Printf("Client sent: %s", req.GetMessage())

  return &buf.PingResponse {
    Reply: "Pong! I received: " + req.GetMessage(),
  }, nil
}

func gRPC() {
  lis, err := net.Listen("tcp", ":50051")
  if err != nil {
    log.Fatalf("Error: %s", err.Error())
  }

  grpcServer := grpc.NewServer()
  buf.RegisterPingServiceServer(grpcServer, &server{})

  log.Println("gRPC Server listening on port 50051")
  if err := grpcServer.Serve(lis); err != nil {
    log.Fatalf("Failed to serve: %s", err.Error())
  }
}
