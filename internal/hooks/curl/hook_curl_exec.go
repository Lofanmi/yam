package curl

import (
	"time"

	"github.com/Lofanmi/yam/api"
)

type curlExec struct {
	request api.PHPRequest
}

func NewExec(request api.PHPRequest) api.HookCallback {
	return &curlExec{request: request}
}

func (s *curlExec) Before(data api.ExecuteData, returnValue api.ZVal) {
	args := data.Args()
	if len(args) != 1 || !args[0].IsResource() {
		return
	}
	id := args[0].AsResource().ID()
	if id <= 0 {
		return
	}
	curl := s.request.Curl(id)
	curl.SetBegin(time.Now())
	curl.BeginTrace(s.request, id)
}

func (s *curlExec) After(data api.ExecuteData, returnValue api.ZVal) {
	args := data.Args()
	if len(args) != 1 || !args[0].IsResource() {
		return
	}
	id := args[0].AsResource().ID()
	if id <= 0 {
		return
	}
	curl := s.request.Curl(id)
	curl.SetEnd(time.Now())
}
