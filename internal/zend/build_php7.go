//go:build php7206 || php7
// +build php7206 php7

package zend

import (
	"unsafe"

	"github.com/Lofanmi/yam/api"
	"github.com/Lofanmi/yam/internal/zend/php7"
)

func NewExecuteData(p unsafe.Pointer) api.ExecuteData {
	return php7.NewExecuteData(p)
}

func NewZVal(p unsafe.Pointer) api.ZVal {
	return php7.NewZVal(p)
}

func NewZArray(p unsafe.Pointer) api.ZArray {
	return php7.NewZArray(p)
}
