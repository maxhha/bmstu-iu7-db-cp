package role

import (
	"auction-back/models"
	"auction-back/ports"
	"fmt"
)

type RolePort struct {
	db ports.DB
}

func New(db ports.DB) RolePort {
	return RolePort{db}
}

func (r *RolePort) HasRole(roleType models.RoleType, viewer models.User) error {
	roles, err := r.db.Role().Find(ports.RoleFindConfig{
		Limit:   1,
		UserIDs: []string{viewer.ID},
		Types:   []models.RoleType{roleType},
	})
	if err != nil {
		return fmt.Errorf("db find: %w", err)
	}

	if len(roles) == 0 {
		return fmt.Errorf("%w: %s", ports.ErrNoRole, roleType)
	}

	return nil
}
