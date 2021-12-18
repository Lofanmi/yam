package api

import (
	"unsafe"
)

type ExecuteData interface {
	Pointer() unsafe.Pointer

	IsClass() bool

	ClassName() string
	FuncName() string
	ClassFuncName() string

	Args() (result []ZVal)

	Object() unsafe.Pointer
}
