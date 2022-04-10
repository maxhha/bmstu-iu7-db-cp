package ports

import "auction-back/models"

//go:generate go run ../codegen/portmocks/main.go --config ../portmocksgen.yml --in role.go --out role_mock.go --outpkg ports

type Role interface {
	HasRole(roleType models.RoleType, viewer models.User) error
}
