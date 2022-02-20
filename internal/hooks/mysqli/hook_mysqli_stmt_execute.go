package mysqli

import (
	"github.com/Lofanmi/yam/api"
)

type mysqliStmtExecute struct {
	request api.PHPRequest
}

func NewStmtExecute(request api.PHPRequest) api.HookCallback {
	return &mysqliStmtExecute{request: request}
}

func (s *mysqliStmtExecute) Before(data api.ExecuteData, returnValue api.ZVal) {
	args := data.Args()
	if len(args) < 1 || !args[0].IsObject() {
		return
	}
	stmtID := uintptr(args[0].AsObject().ID())

	m := s.request.MySQLiStmtMapping()
	id := m[stmtID]

	r := s.request.MySQLi(id)
	r.BeginTrace(s.request, id)
}

func (s *mysqliStmtExecute) After(data api.ExecuteData, returnValue api.ZVal) {
	args := data.Args()
	if len(args) < 1 || !args[0].IsObject() {
		return
	}
	stmtID := uintptr(args[0].AsObject().ID())

	m := s.request.MySQLiStmtMapping()
	id := m[stmtID]

	r := s.request.MySQLi(id)
	r.EndTrace(s.request, id)
}
