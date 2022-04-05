package token_sender

import (
	"auction-back/grpc/notifier"
	"auction-back/models"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ReceiverGetter func(token models.Token) (string, error)
type DataGetter func(token models.Token) (map[string]string, error)

type Config struct {
	Name              string
	AddressEnvVarName string
	ReceiverGetters   map[models.TokenAction]ReceiverGetter
	DataGetters       map[models.TokenAction]DataGetter
}

type TokenSender struct {
	client notifier.NotifierClient
	config Config
}

func New(config Config) *TokenSender {
	address, ok := os.LookupEnv(config.AddressEnvVarName)
	if !ok {
		log.Fatalf("%s does not exist in environment variables!", config.AddressEnvVarName)
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("email notifier grpc dial: %v", err)
	}

	return &TokenSender{
		client: notifier.NewNotifierClient(conn),
		config: config,
	}
}

func (n *TokenSender) Name() string {
	return n.config.Name
}

func (n *TokenSender) Send(token models.Token) (bool, error) {
	action := token.Action

	getReceiver, shouldSend := n.config.ReceiverGetters[action]
	if !shouldSend {
		return false, nil
	}

	receiver, err := getReceiver(token)
	if err != nil {
		return false, fmt.Errorf("get receiver for %s: %w", token.Action, err)
	}

	var data map[string]string
	if getData, has := n.config.DataGetters[token.Action]; has {
		var err error
		if data, err = getData(token); err != nil {
			return false, fmt.Errorf("get data for %s: %w", token.Action, err)
		}
	}

	input := notifier.SendInput{
		Receivers: []string{receiver},
		Action:    string(token.Action),
		Data:      data,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := n.client.Send(ctx, &input)
	if err != nil {
		return false, fmt.Errorf("client: %w", err)
	}

	if result.Status != "OK" {
		return false, fmt.Errorf("client status: %v", result.Status)
	}

	return true, nil
}
