package mysqli

import (
	"github.com/Lofanmi/yam/api"
	"github.com/spf13/cast"
)

type mysqliConnect struct {
	request api.PHPRequest
}

func NewConnect(request api.PHPRequest) api.HookCallback {
	return &mysqliConnect{request: request}
}

func (s *mysqliConnect) Before(data api.ExecuteData, returnValue api.ZVal) {}

func (s *mysqliConnect) After(data api.ExecuteData, returnValue api.ZVal) {
	if returnValue.IsNull() {
		return
	}
	args := data.Args()
	if len(args) < 4 {
		return
	}
	r := s.request.MySQLi(uintptr(returnValue.AsObject().ID()))
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
