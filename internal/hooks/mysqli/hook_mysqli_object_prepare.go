package mysqli

import (
	"strings"

	"github.com/Lofanmi/yam/api"
	"github.com/Lofanmi/yam/internal/zend"
)

type mysqliObjectPrepare struct {
	request api.PHPRequest
}

func NewObjectPrepare(request api.PHPRequest) api.HookCallback {
	return &mysqliObjectPrepare{request: request}
}

func (s *mysqliObjectPrepare) Before(data api.ExecuteData, returnValue api.ZVal) {}

func (s *mysqliObjectPrepare) After(data api.ExecuteData, returnValue api.ZVal) {
	args := data.Args()
	if len(args) < 1 || !args[0].IsString() {
		return
	}

	id := uintptr(zend.ExecuteDataObjectID(data))
	r := s.request.MySQLi(id)

	sql := args[0].AsString()
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
