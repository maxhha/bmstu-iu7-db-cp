package test

import (
	"database/sql/driver"
	"reflect"
	"sync"

	"github.com/DATA-DOG/go-sqlmock"
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
