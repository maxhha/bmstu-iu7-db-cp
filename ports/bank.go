package ports

import "auction-back/models"

//go:generate go run ../codegen/portmocks/main.go --config ../portmocksgen.yml --in bank.go --out bank_mock.go --outpkg ports

type Bank interface {
	UserFormApproved(form models.UserForm) error
}
