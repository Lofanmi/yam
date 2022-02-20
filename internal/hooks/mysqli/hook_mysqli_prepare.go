package mysqli

import (
	"strings"

	"github.com/Lofanmi/yam/api"
)

type mysqliPrepare struct {
	request api.PHPRequest
}

func NewPrepare(request api.PHPRequest) api.HookCallback {
	return &mysqliPrepare{request: request}
}

func (s *mysqliPrepare) Before(data api.ExecuteData, returnValue api.ZVal) {}

func (s *mysqliPrepare) After(data api.ExecuteData, returnValue api.ZVal) {
	args := data.Args()
	if len(args) < 2 || !args[0].IsObject() || !args[1].IsString() {
		return
	}

	id := uintptr(args[0].AsObject().ID())
	r := s.request.MySQLi(id)

	sql := strings.TrimSpace(args[1].AsString())
	pieces := strings.SplitN(sql, " ", 2)
	if len(pieces) < 1 || len(pieces[0]) <= 0 {
		return
	}

	r.SetOperation(strings.ToLower(pieces[0]))
	r.SetStatement(sql)

	if returnValue.IsObject() {
		stmtID := uintptr(returnValue.AsObject().ID())
		m := s.request.MySQLiStmtMapping()
		m[stmtID] = id
	}
}
