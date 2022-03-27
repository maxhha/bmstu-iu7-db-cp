package test

import (
	"database/sql/driver"
	"reflect"
	"sync"

	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func MockRows(objs ...interface{}) *sqlmock.Rows {
	s, err := schema.Parse(objs[0], &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		panic("failed to create schema")
	}

	columns := make([]string, 0)
	for _, field := range s.Fields {
		if len(field.DBName) == 0 {
			continue
		}
		columns = append(columns, field.DBName)
	}

	rows := sqlmock.NewRows(columns)

	for _, obj := range objs {
		row := make([]driver.Value, 0)

		for _, field := range s.Fields {
			if len(field.DBName) == 0 {
				continue
			}
			r := reflect.ValueOf(obj)
			f := reflect.Indirect(r).FieldByName(field.Name)
			row = append(row, f.Interface())
		}

		rows = rows.AddRow(row...)
	}

	return rows
}

type DBSuite struct {
	suite.Suite
	SqlDB   *sql.DB
	DB      *gorm.DB
	SqlMock sqlmock.Sqlmock
}

func (s *DBSuite) SetupTest() {
	var err error
	s.SqlDB, s.SqlMock, err = sqlmock.New()
	require.NoError(s.T(), err)
	require.NotNil(s.T(), s.SqlDB)
	require.NotNil(s.T(), s.SqlMock)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 s.SqlDB,
		PreferSimpleProtocol: true,
	})

	s.DB, err = gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
	})
	require.NoError(s.T(), err)
}

func (s *DBSuite) TearDownTest() {
	s.SqlDB.Close()
}
