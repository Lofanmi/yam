package redis

import (
	"strconv"
	"strings"
	"time"

	"github.com/Lofanmi/yam/api"
	"github.com/spf13/cast"
)

type phpRedis struct {
	request api.PHPRequest
}

func NewPHPRedis(request api.PHPRequest) api.HookCallback {
	return &phpRedis{request: request}
}

func (s *phpRedis) Before(data api.ExecuteData, returnValue api.ZVal) {
	fn := data.FuncName()
	if fn == "__construct" || fn == "__destruct" {
		return
	}
	ptr := data.Object()
	operation := strings.ToLower(fn)
	r := s.request.PHPRedis(uintptr(ptr))
	r.SetBegin(time.Now())
	r.SetOperation(operation)
	if operation == "auth" {
		r.SetStatement("auth ********")
	} else {
		var sb strings.Builder
		sb.WriteString(operation)
		for _, val := range data.Args() {
			sb.WriteByte(' ')
			sb.WriteString(cast.ToString(val.Value()))
		}
		r.SetStatement(sb.String())
	}
	if operation == "connect" {
		pieces := strings.Split(r.Statement(), " ")
		if len(pieces) >= 3 {
			r.SetHost(pieces[1])
			r.SetPort(pieces[2])
		}
	}
	r.BeginTrace(s.request, uintptr(ptr))
}

func (s *phpRedis) After(data api.ExecuteData, returnValue api.ZVal) {
	fn := strings.ToLower(data.FuncName())
	if fn == "__construct" || fn == "__destruct" {
		return
	}
	ptr := data.Object()
	r := s.request.PHPRedis(uintptr(ptr))
	if r.Host() == "" || r.Port() == "" {
		host, port, ok := api.RedisAddr(ptr)
		if !ok {
			return
		}
		r.SetHost(host)
		r.SetPort(strconv.FormatInt(int64(port), 10))
	}
	r.SetEnd(time.Now())
	r.EndTrace(s.request, uintptr(ptr))
}
