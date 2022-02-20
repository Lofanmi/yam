package php7

import (
	"strconv"
	"strings"
	"unsafe"

	"github.com/Lofanmi/yam/api"
)

var _ api.ExecuteData = (*executeData)(nil)

const (
	unknownFuncName = "[UnknownFuncName]"
)

type executeData struct {
	p unsafe.Pointer
}

func NewExecuteData(p unsafe.Pointer) api.ExecuteData {
	return &executeData{p}
}

func (e *executeData) Pointer() unsafe.Pointer {
	return e.p
}

func (e *executeData) IsClass() bool {
	p := (*zendExecuteData)(e.p)
	if p == nil || p.Func == nil || p.Func.Common.Scope == nil {
		return false
	}
	if p.Func.Common.Scope.Name == nil {
		return false
	}
	return p.Func.Common.Scope.Name.String() != ""
}

func (e *executeData) ClassName() string {
	if !e.IsClass() {
		return ""
	}
	p := (*zendExecuteData)(e.p)
	return p.Func.Common.Scope.Name.String()
}

func (e *executeData) FuncName() string {
	p := (*zendExecuteData)(e.p)
	if p == nil || p.Func == nil {
		return unknownFuncName
	}
	return p.Func.Common.FunctionName.String()
}

func (e *executeData) ClassFuncName() string {
	if e.IsClass() {
		return e.ClassName() + "::" + e.FuncName() + "()"
	}
	return e.FuncName() + "()"
}

func (e *executeData) Args() (result []api.ZVal) {
	p := (*zendExecuteData)(e.p)
	argumentCount := int(p.This.U2)
	for i := 0; i < argumentCount; i++ {
		v := unsafe.Pointer(uintptr(e.p) + unsafe.Sizeof(zendExecuteData{}) + uintptr(i)*unsafe.Sizeof(zVal{}))
		if v == nil {
			continue
		}
		result = append(result, NewZVal(v))
	}
	return
}

func (e *executeData) RedisHostPort() (host string, port int) {
	p := (*redisObject)((*zendExecuteData)(e.p).This.AsPointer())
	host, port = p.Sock.Host.String(), int(p.Sock.Port)
	return
}

func (e *executeData) PdoHostPort() (host string, port int) {
	p := (*pdoDbh)((*pdoDbhObject)((*zendExecuteData)(e.p).This.AsPointer()).inner)
	dsnBytes := make([]byte, 0, p.dataSourceLen)
	var k unsafe.Pointer
	for i := 0; i < p.dataSourceLen; i++ {
		k = unsafe.Pointer(uintptr(p.dataSource) + uintptr(i))
		dsnBytes[i] = *(*byte)(k)
	}
	for _, kv := range strings.Split(string(dsnBytes), ";") {
		pieces := strings.SplitN(kv, "=", 2)
		if len(pieces) == 2 && strings.ToLower(pieces[0]) == "host" {
			host = pieces[1]
			continue
		}
		if len(pieces) == 2 && strings.ToLower(pieces[0]) == "port" {
			if _port, _err := strconv.ParseInt(pieces[1], 10, 64); _err == nil {
				port = int(_port)
				return
			}
		}
	}
	return
}
