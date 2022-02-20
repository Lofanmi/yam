package mysqli

import (
	"strings"

	"github.com/Lofanmi/yam/api"
	"github.com/Lofanmi/yam/internal/zend"
)

type mysqliObjectQuery struct {
	request api.PHPRequest
}

func NewObjectQuery(request api.PHPRequest) api.HookCallback {
	return &mysqliObjectQuery{request: request}
}

func (s *mysqliObjectQuery) Before(data api.ExecuteData, returnValue api.ZVal) {
	args := data.Args()
	if len(args) < 1 || !args[0].IsString() {
		return
	}

	id := uintptr(zend.ExecuteDataObjectID(data))
	r := s.request.MySQLi(id)

	sql := strings.TrimSpace(args[0].AsString())
	pieces := strings.SplitN(sql, " ", 2)
	if len(pieces) < 1 || len(pieces[0]) <= 0 {
		return
	}

	r.SetOperation(strings.ToLower(pieces[0]))
	r.SetStatement(sql)

	r.BeginTrace(s.request, id)
}

func (s *mysqliObjectQuery) After(data api.ExecuteData, returnValue api.ZVal) {
	args := data.Args()
	if len(args) < 1 || !args[0].IsString() {
		return
	}

	id := uintptr(zend.ExecuteDataObjectID(data))
	r := s.request.MySQLi(id)

	r.EndTrace(s.request, id)
}
