package notifier

import (
	"auction-back/grpc/notifier"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func New(addressEnv string) notifier.NotifierClient {
	address, ok := os.LookupEnv(addressEnv)
	if !ok {
		log.Fatalf("%s does not exist in environment variables!", addressEnv)
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("notifer dial: %v", err)
	}

	return notifier.NewNotifierClient(conn)
}
