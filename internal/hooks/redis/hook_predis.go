package redis

import (
	"net/url"
	"strings"
	"time"

	"github.com/Lofanmi/yam/api"
	"github.com/spf13/cast"
)

type pRedis struct {
	request api.PHPRequest
}

func NewPRedis(request api.PHPRequest) api.HookCallback {
	return &pRedis{request: request}
}

func parseAddr(s string) (host, port string) {
	if strings.Contains(s, "unix") {
		host, port = s, ""
		return
	}
	u, err := url.Parse(s)
	if err != nil {
		host, port = s, ""
	}
	host, port = u.Hostname(), u.Port()
	return
}

func (s *pRedis) Before(data api.ExecuteData, returnValue api.ZVal) {
	fn := data.FuncName()
	if !(fn == "__construct" || fn == "__call") {
		return
	}
	if fn == "__construct" {
		args := data.Args()
		ptr := data.Object()
		r := s.request.PHPRedis(uintptr(ptr))
		host, port := "(unknown host)", "(unknown port)"
		if len(args) <= 0 {
			host, port = "127.0.0.1", "6379"
		} else if args[0].IsString() {
			host, port = parseAddr(args[0].AsString())
		} else if args[0].IsArray() {
			iz := args[0].AsArray().ToIZMap()
			if addr, exist := iz[0]; exist && addr.IsString() {
				host, port = parseAddr(addr.AsString())
			} else if scheme, ok := iz["scheme"]; ok && scheme.IsString() {
				switch scheme.AsString() {
				case "tcp":
					if _host, ok2 := iz["host"]; ok2 && _host.IsString() {
						host = _host.AsString()
					}
					if _port, ok2 := iz["port"]; ok2 {
						port = cast.ToString(_port.Value())
					}
				case "unix":
					if path, ok2 := iz["path"]; ok2 && path.IsString() {
						host, port = "unix:"+path.AsString(), ""
					}
				case "tls":
					if iz2, ok2 := iz["ssl"]; ok2 && iz2.IsArray() {
						params := iz2.AsArray().ToIZMap()
						if f, ok3 := params["cafile"]; ok3 && f.IsString() {
							host, port = f.AsString(), ""
						}
					}
				}
			}
		}
		r.SetHost(host)
		r.SetPort(port)
		return
	}
	args := data.Args()
	if len(args) != 2 {
		return
	}
	var (
		sb        strings.Builder
		operation string
	)
	ptr := data.Object()
	r := s.request.PHPRedis(uintptr(ptr))
	r.SetBegin(time.Now())
	if args[0].IsString() {
		operation = args[0].AsString()
		sb.WriteString(operation)
		r.SetOperation(operation)
		if operation == "auth" {
			r.SetStatement("auth ********")
		}
	}
	if args[1].IsArray() && operation != "auth" {
		iz := args[1].AsArray().ToIZMap()
		for i := 0; i < len(iz); i++ {
			if val, ok := iz[i]; ok {
				sb.WriteByte(' ')
				sb.WriteString(cast.ToString(val.Value()))
			}
		}
		r.SetStatement(sb.String())
	}
	r.BeginTrace(s.request, uintptr(ptr))
}

func (s *pRedis) After(data api.ExecuteData, returnValue api.ZVal) {
	fn := strings.ToLower(data.FuncName())
	if fn != "__call" {
		return
	}
	ptr := data.Object()
	r := s.request.PHPRedis(uintptr(ptr))
	r.SetEnd(time.Now())
	r.EndTrace(s.request, uintptr(ptr))
}
