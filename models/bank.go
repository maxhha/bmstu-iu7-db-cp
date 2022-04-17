package models

import (
	"database/sql"
	"time"
)

type Bank struct {
	ID                   string
	Name                 string
	Bic                  string
	CorrespondentAccount string
	Inn                  string
	Kpp                  string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            sql.NullTime
}
