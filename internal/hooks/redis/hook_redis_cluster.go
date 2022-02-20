package redis

import (
	"strings"
	"time"

	"github.com/Lofanmi/yam/api"
	"github.com/spf13/cast"
)

type phpRedisCluster struct {
	request api.PHPRequest
}

func NewPHPRedisCluster(request api.PHPRequest) api.HookCallback {
	return &phpRedisCluster{request: request}
}

func (s *phpRedisCluster) Before(data api.ExecuteData, returnValue api.ZVal) {
	fn := data.FuncName()
	if fn == "__construct" {
		args := data.Args()
		if len(args) >= 2 && args[1].IsArray() {
			iz := args[1].AsArray().ToIZMap()
			if addr, ok := iz[0]; ok && addr.IsString() {
				pieces := strings.SplitN(addr.AsString(), ":", 2)
				if len(pieces) == 2 {
					ptr := data.Object()
					r := s.request.PHPRedis(uintptr(ptr))
					r.SetHost(pieces[0])
					r.SetPort(pieces[1])
				}
			}
		}
		return
	}
	if fn == "__destruct" {
		return
	}
	ptr := data.Object()
	operation := strings.ToLower(fn)
	r := s.request.PHPRedis(uintptr(ptr))
	r.SetBegin(time.Now())
	r.SetOperation(operation)
	var sb strings.Builder
	sb.WriteString(operation)
	for _, val := range data.Args() {
		sb.WriteByte(' ')
		sb.WriteString(cast.ToString(val.Value()))
	}
	r.SetStatement(sb.String())
	r.BeginTrace(s.request, uintptr(ptr))
}

func (s *phpRedisCluster) After(data api.ExecuteData, returnValue api.ZVal) {
	fn := strings.ToLower(data.FuncName())
	if fn == "__construct" || fn == "__destruct" {
		return
	}
	ptr := data.Object()
	r := s.request.PHPRedis(uintptr(ptr))
	r.SetEnd(time.Now())
	r.EndTrace(s.request, uintptr(ptr))
}
