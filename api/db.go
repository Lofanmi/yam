package api

import (
	"time"
)

type DB interface {
	Begin() (t time.Time)
	End() (t time.Time)

	Host() string
	Port() string
	Addr() string
	Operation() string
	Statement() string

	User() string
	Database() string

	SetBegin(t time.Time)
	SetEnd(t time.Time)
	SetHost(v string)
	SetPort(v string)
	SetOperation(v string)
	SetStatement(v string)

	SetUser(v string)
	SetDatabase(v string)

	BeginTrace(request PHPRequest, id uintptr)
	EndTrace(request PHPRequest, id uintptr)
}
