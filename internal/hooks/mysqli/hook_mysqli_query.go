package mysqli

import (
	"strings"

	"github.com/Lofanmi/yam/api"
)

type mysqliQuery struct {
	request api.PHPRequest
}

func NewQuery(request api.PHPRequest) api.HookCallback {
	return &mysqliQuery{request: request}
}

func (s *mysqliQuery) Before(data api.ExecuteData, returnValue api.ZVal) {
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

	r.BeginTrace(s.request, id)
}

func (s *mysqliQuery) After(data api.ExecuteData, returnValue api.ZVal) {
	args := data.Args()
	if len(args) < 2 || !args[0].IsObject() || !args[1].IsString() {
		return
	}

	id := uintptr(args[0].AsObject().ID())
	r := s.request.MySQLi(id)

	r.EndTrace(s.request, id)
}
