package main

import (
	"auction-back/grpc/notifier"
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

type NotifierServer struct {
	notifier.UnimplementedNotifierServer
}

func New() *NotifierServer {
	return &NotifierServer{}
}

func (s *NotifierServer) Send(ctx context.Context, input *notifier.SendInput) (*notifier.SendResult, error) {
	fmt.Printf("\nreceivers: %v\naction: %v\ndata: %v\n", input.Receivers, input.Action, input.Data)
	// TODO: send to redis
	return &notifier.SendResult{Status: "OK"}, nil
}

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatalln("PORT is not defiend in environment variables!")
	}

	host, _ := os.LookupEnv("HOST")
	address := fmt.Sprintf("%s:%s", host, port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalln(err)
	}

	grpcServer := grpc.NewServer()
	notifier.RegisterNotifierServer(grpcServer, New())

	fmt.Printf("Start listen on %s\n", address)
	grpcServer.Serve(lis)
}
