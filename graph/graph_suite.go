package graph

import (
	"auction-back/ports"

	"github.com/stretchr/testify/suite"
)

type GraphSuite struct {
	suite.Suite
	DB        ports.DBMock
	TokenMock ports.TokenMock
	BankMock  ports.BankMock
	RoleMock  ports.RoleMock
	resolver  *Resolver
}

func (s *GraphSuite) SetupTest() {
	s.resolver = New(&s.DB, &s.TokenMock, &s.BankMock, &s.RoleMock, nil)
}
