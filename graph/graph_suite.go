package graph

import (
	"auction-back/test"

	"github.com/stretchr/testify/suite"
)

type GraphSuite struct {
	suite.Suite
	DB        test.DBMock
	TokenMock test.TokenPort
	BankMock  test.BankPort
	RoleMock  test.RolePort
	resolver  *Resolver
}

func (s *GraphSuite) SetupTest() {
	s.resolver = New(&s.DB, &s.TokenMock, &s.BankMock, &s.RoleMock)
}
