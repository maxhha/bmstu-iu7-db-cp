package test

import (
	"database/sql/driver"
	"time"
)

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type AfterTime struct {
	time.Time
}

func (a AfterTime) Match(v driver.Value) bool {
	t, ok := v.(time.Time)

	if !ok {
		return false
	}

	return t.After(a.Time)
}
