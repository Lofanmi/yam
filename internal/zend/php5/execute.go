package php5

import "C"
import (
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
	if p == nil || p.FunctionState.Function == nil || p.FunctionState.Function.Common.Scope == nil {
		return false
	}
	if p.FunctionState.Function.Common.Scope.Name == nil {
		return false
	}
	return C.GoString((*C.char)(p.FunctionState.Function.Common.Scope.Name)) != ""
}

func (e *executeData) ClassName() string {
	if !e.IsClass() {
		return ""
	}
	p := (*zendExecuteData)(e.p)
	return C.GoString((*C.char)(p.FunctionState.Function.Common.Scope.Name))
}

func (e *executeData) FuncName() string {
	p := (*zendExecuteData)(e.p)
	if p == nil || p.FunctionState.Function == nil || p.FunctionState.Function.Common.FunctionName == nil {
		return unknownFuncName
	}
	return C.GoString((*C.char)(p.FunctionState.Function.Common.FunctionName))
}

func (e *executeData) ClassFuncName() string {
	if e.IsClass() {
		return e.ClassName() + "::" + e.FuncName() + "()"
	}
	return e.FuncName() + "()"
}

func (e *executeData) Args() (result []api.ZVal) {
	p := (*zendExecuteData)(e.p)
	a := p.FunctionState.Arguments
	argumentCount := int(*(*uint32)(a))
	for i := 0; i < argumentCount; i++ {
		v := unsafe.Pointer(uintptr(a) - uintptr(argumentCount-i)*unsafe.Sizeof(0))
		if v == nil {
			continue
		}
		result = append(result, NewZVal(v))
	}
	return
}

func (e *executeData) Object() unsafe.Pointer {
	p := (*zendExecuteData)(e.p)
	return unsafe.Pointer(p.Object)
}
