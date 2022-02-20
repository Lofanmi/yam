package api

import (
	"unsafe"
)

type ZVal interface {
	Value() interface{}

	IsInt() bool
	IsNull() bool
	IsBool() bool
	IsFloat() bool
	IsString() bool
	IsArray() bool
	IsResource() bool
	IsPointer() bool
	IsObject() bool

	AsInt() int
	AsBool() bool
	AsFloat() float64
	AsString() string
	AsArray() ZArray
	AsResource() ZResource
	AsPointer() unsafe.Pointer
	AsObject() ZObject
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

type ZObject interface {
	ID() int32
	Pointer() unsafe.Pointer
}
