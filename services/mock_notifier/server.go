package main

import (
	"auction-back/grpc/notifier"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
)

type NotifierServer struct {
	notifier.UnimplementedNotifierServer
	rdb *redis.Client
}

func New(rdb *redis.Client) *NotifierServer {
	return &NotifierServer{rdb: rdb}
}

func (s *NotifierServer) Send(ctx context.Context, input *notifier.SendInput) (*notifier.SendResult, error) {
	key := fmt.Sprintf("send-%s-%s", input.Action, strings.Join(input.Receivers, ","))
	fmt.Printf("\n%s\nreceivers: %v\naction: %v\ndata: %v\n", key, input.Receivers, input.Action, input.Data)

	data, err := json.Marshal(input.Data)
	if err != nil {
		return &notifier.SendResult{Status: fmt.Sprintf("json: %v", err)}, nil
	}

	if err := s.rdb.Set(ctx, key, data, 0).Err(); err != nil {
		return &notifier.SendResult{Status: fmt.Sprintf("rdb set: %v", err)}, nil
	}

	return &notifier.SendResult{Status: "OK"}, nil
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatalln("PORT is not defiend in environment variables!")
	}

	host := os.Getenv("HOST")
	address := fmt.Sprintf("%s:%s", host, port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalln(err)
	}

	grpcServer := grpc.NewServer()
	notifier.RegisterNotifierServer(grpcServer, New(rdb))

	fmt.Printf("Start listen on %s\n", address)
	grpcServer.Serve(lis)
}
