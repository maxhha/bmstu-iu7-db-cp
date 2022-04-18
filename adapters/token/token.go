package token

import (
	"auction-back/adapters/token/mock"
	"auction-back/adapters/token/prod"
	"auction-back/ports"
	"log"
	"os"
)

var creators = map[string]func(db ports.DB) ports.Token{
	"PROD": func(db ports.DB) ports.Token {
		senders := []ports.TokenSender{
			emailTokenSender(db),
			phoneTokenSender(db),
		}

		tokenAdapter := prod.New(db, senders)
		return &tokenAdapter
	},
	"MOCK": func(db ports.DB) ports.Token {
		tokenAdapter := mock.New(db)
		return &tokenAdapter
	},
}

func New(db ports.DB) ports.Token {
	tokenAdapterName, ok := os.LookupEnv("TOKEN_ADAPTER")
	if !ok || tokenAdapterName == "" {
		tokenAdapterName = "PROD"
	}

	creator, exists := creators[tokenAdapterName]
	if !exists {
		log.Fatalf("Unknown TOKEN_ADAPTER = '%s'\n", tokenAdapterName)
	}

	return creator(db)
}
