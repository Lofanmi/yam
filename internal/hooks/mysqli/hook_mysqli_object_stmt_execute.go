package mysqli

import (
	"github.com/Lofanmi/yam/api"
	"github.com/Lofanmi/yam/internal/zend"
)

type mysqliObjectStmtExecute struct {
	request api.PHPRequest
}

func NewObjectStmtExecute(request api.PHPRequest) api.HookCallback {
	return &mysqliObjectStmtExecute{request: request}
}

func (s *mysqliObjectStmtExecute) Before(data api.ExecuteData, returnValue api.ZVal) {
	stmtID := uintptr(zend.ExecuteDataObjectID(data))

	m := s.request.MySQLiStmtMapping()
	id := m[stmtID]

	r := s.request.MySQLi(id)
	r.BeginTrace(s.request, id)
}

func (s *mysqliObjectStmtExecute) After(data api.ExecuteData, returnValue api.ZVal) {
	stmtID := uintptr(zend.ExecuteDataObjectID(data))

	m := s.request.MySQLiStmtMapping()
	id := m[stmtID]

	r := s.request.MySQLi(id)
	r.EndTrace(s.request, id)
}
