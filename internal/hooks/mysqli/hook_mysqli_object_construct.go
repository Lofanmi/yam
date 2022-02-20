package mysqli

import (
	"github.com/Lofanmi/yam/api"
	"github.com/Lofanmi/yam/internal/zend"
	"github.com/spf13/cast"
)

type mysqliObjectConstruct struct {
	request api.PHPRequest
}

func NewObjectConstruct(request api.PHPRequest) api.HookCallback {
	return &mysqliObjectConstruct{request: request}
}

func (s *mysqliObjectConstruct) Before(data api.ExecuteData, returnValue api.ZVal) {}

func (s *mysqliObjectConstruct) After(data api.ExecuteData, returnValue api.ZVal) {
	args := data.Args()
	if len(args) < 4 {
		return
	}
	id := uintptr(zend.ExecuteDataObjectID(data))
	r := s.request.MySQLi(id)
	r.SetPort("3306")
	if args[0].IsString() {
		r.SetHost(args[0].AsString())
	}
	if args[1].IsString() {
		r.SetUser(args[1].AsString())
	}
	if args[3].IsString() {
		r.SetDatabase(args[3].AsString())
	}
	if len(args) >= 5 {
		r.SetPort(cast.ToString(args[4].Value()))
	}
}
