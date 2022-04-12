package ports

import "auction-back/models"

//go:generate go run ../codegen/portmocks/main.go --config ../portmocksgen.yml --in token.go --out token_mock.go --outpkg ports

type Token interface {
	Create(action models.TokenAction, viewer models.User, data map[string]interface{}) error
	Activate(action models.TokenAction, token_code string, viewer models.User) (models.Token, error)
}

type TokenSender interface {
	Name() string
	Send(models.Token) (bool, error)
}
