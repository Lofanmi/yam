package api

import (
	"unsafe"
)

type ZVal interface {
	Value() interface{}

	IsInt() bool
	IsBool() bool
	IsFloat() bool
	IsString() bool
	IsArray() bool
	IsResource() bool
	IsPointer() bool

	AsInt() int
	AsBool() bool
	AsFloat() float64
	AsString() string
	AsArray() ZArray
	AsResource() ZResource
	AsPointer() unsafe.Pointer
}

type ZArray interface {
	Len() int
	ToSSMap() (m map[string]string)
	ToSZMap() (m map[string]ZVal)
	ToIIMap() (m map[interface{}]interface{})
	ToIZMap() (m map[interface{}]ZVal)
}

type ZResource interface {
	ID() int32
	Pointer() unsafe.Pointer
}
