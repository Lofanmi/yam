//go:build php5
// +build php5

package zend

import (
	"unsafe"

	"github.com/Lofanmi/yam/api"
	"github.com/Lofanmi/yam/internal/zend/php5"
)

type String = php5.ZendString

func NewExecuteData(p unsafe.Pointer) api.ExecuteData {
	return php5.NewExecuteData(p)
}

func NewZVal(p unsafe.Pointer) api.ZVal {
	return php5.NewZVal(p)
}

func NewZArray(p unsafe.Pointer) api.ZArray {
	return php5.NewZArray(p)
}

func ExecuteDataObjectID(data api.ExecuteData) int32 {
	p := data.Object()
	z := NewZVal(unsafe.Pointer(&p))
	if z.IsObject() {
		return z.AsObject().ID()
	}
	return 0
}
