package php5

import "C"
import (
	"fmt"
	"unsafe"

	"github.com/Lofanmi/yam/api"
)

var _ api.ZVal = (*zVal)(nil)

type (
	ZType      byte
	ZendString = zendString
)

const (
	IsNull          ZType = 0
	IsLong          ZType = 1
	IsDouble        ZType = 2
	IsBool          ZType = 3
	IsArray         ZType = 4
	IsObject        ZType = 5
	IsString        ZType = 6
	IsResource      ZType = 7
	IsConstant      ZType = 8
	IsConstantArray ZType = 9
	IsCallable      ZType = 10

	ZendSuccess = 0
	ZendFailure = -1

	HashKeyIsString = 1 // HASH_KEY_IS_STRING
	HashKeyIsLong   = 2 // HASH_KEY_IS_LONG
	HashKeyNonExist = 3 // HASH_KEY_NON_EXISTANT

	HashUpdate     = 1 << 0 // HASH_UPDATE
	HashAdd        = 1 << 1 // HASH_ADD
	HashNextInsert = 1 << 2 // HASH_NEXT_INSERT

	HashDelKey      = 0 // HASH_DEL_KEY
	HashDelIndex    = 1 // HASH_DEL_INDEX
	HashDelKeyQuick = 2 // HASH_DEL_KEY_QUICK

	HashUpdateKeyIfNone   = 0 // HASH_UPDATE_KEY_IF_NONE
	HashUpdateKeyIfBefore = 1 // HASH_UPDATE_KEY_IF_BEFORE
	HashUpdateKeyIfAfter  = 2 // HASH_UPDATE_KEY_IF_AFTER
	HashUpdateKeyAnyway   = 3 // HASH_UPDATE_KEY_ANYWAY
)

type (
	zVal struct {
		ZendValue  zValueValue
		RefCountGC uint32
		ZendType   byte
		IsRefGC    uint8
	}
	zValueValue struct {
		A, B unsafe.Pointer
	}
	zendClassEntry struct {
		Type byte
		Name unsafe.Pointer
		// ...
	}
	zendFunction struct {
		// Type   uint8 /* MUST be the first element of this struct! */
		Common struct {
			Type            uint8 /* never used */
			FunctionName    unsafe.Pointer
			Scope           *zendClassEntry // zend_class_entry
			FnFlags         uint32
			Prototype       unsafe.Pointer // zend_function
			NumArgs         uint32
			RequiredNumArgs uint32
			ArgInfo         unsafe.Pointer // zend_arg_info
		}
		// OpArray
		// InternalFunction
	}
	zendFunctionState struct {
		Function  *zendFunction
		Arguments unsafe.Pointer
	}
	zendExecuteData struct {
		OpLine              unsafe.Pointer
		FunctionState       zendFunctionState
		FBC                 *zendFunction
		CalledScope         *zendClassEntry
		OpArray             unsafe.Pointer
		Object              *zVal
		TS                  unsafe.Pointer
		CVS                 ***zVal
		SymbolTable         *hashTable
		PrevExecuteData     *zendExecuteData
		OldErrorReporting   *zVal
		Nested              uint8
		OriginalReturnValue **zVal
		CurrentScope        *zendClassEntry
		CurrentCalledScope  *zendClassEntry
		CurrentThis         *zVal
		CurrentObject       *zVal
	}
	hashTable struct {
		TableSize       uint32
		TableMask       uint32
		NumOfElements   uint32
		NextFreeElement uint64
		InternalPointer *bucket /* Used for element traversal */
		ListHead        *bucket
		ListTail        *bucket
		ArBuckets       **bucket
		Destructor      unsafe.Pointer
		Persistent      uint8
		ApplyCount      uint8
		ApplyProtection uint8
	}
	bucket struct {
		H         uint64
		KeyLength uint32
		Data      unsafe.Pointer
		DataPtr   unsafe.Pointer
		ListNext  *bucket
		ListLast  *bucket
		Next      *bucket
		Last      *bucket
		ArKey     unsafe.Pointer
	}
	zendResource struct {
		val *zVal
	}
	zendObject struct {
		val *zVal
	}
	zendString struct {
		Gc  int16
		Len int
		Val unsafe.Pointer
	}
)

func NewZVal(p unsafe.Pointer) api.ZVal {
	if p == nil {
		return nil
	}
	from := new(*zVal)
	from = (**zVal)(p)
	return &zVal{
		ZendValue: zValueValue{
			A: (*from).ZendValue.A,
			B: (*from).ZendValue.B,
		},
		RefCountGC: (*from).RefCountGC,
		ZendType:   (*from).ZendType,
		IsRefGC:    (*from).IsRefGC,
	}
}

func (z *zVal) Type() int {
	if z == nil {
		return 0
	}
	return int(z.ZendType)
}

func (z *zVal) Value() interface{} {
	if z == nil {
		return nil
	}
	switch ZType(z.ZendType) {
	case IsNull:
		return nil
	case IsBool:
		return !(uint64(uintptr(z.ZendValue.A)) == 0)
	case IsLong:
		return z.AsInt()
	case IsDouble:
		return z.AsFloat()
	case IsString:
		return z.AsString()
	case IsArray:
		return z.AsArray()
	case IsResource:
		return z.AsResource()
	case IsObject:
		return z.AsObject()
	default:
		return fmt.Sprintf("unknown type: %d %p %p", ZType(z.ZendType), z.ZendValue.A, z.ZendValue.B)
	}
}

func (z *zVal) IsNull() bool {
	return ZType(z.ZendType) == IsNull
}

func (z *zVal) IsBool() bool {
	return ZType(z.ZendType) == IsBool
}

func (z *zVal) IsInt() bool {
	return ZType(z.ZendType) == IsLong
}

func (z *zVal) IsString() bool {
	return ZType(z.ZendType) == IsString
}

func (z *zVal) IsFloat() bool {
	return ZType(z.ZendType) == IsDouble
}

func (z *zVal) IsArray() bool {
	return ZType(z.ZendType) == IsArray
}

func (z *zVal) IsResource() bool {
	return ZType(z.ZendType) == IsResource
}

func (z *zVal) IsPointer() bool {
	return false
}

func (z *zVal) IsObject() bool {
	return ZType(z.ZendType) == IsObject
}

func (z *zVal) AsInt() int {
	if z == nil {
		return 0
	}
	return int(uintptr(z.ZendValue.A))
}

func (z *zVal) AsBool() bool {
	if z == nil {
		return false
	}
	return !(uint64(uintptr(z.ZendValue.A)) == 0)
}

func (z *zVal) AsFloat() float64 {
	if z == nil {
		return 0.0
	}
	p := new(float64)
	i := (*int)(unsafe.Pointer(p))
	*i = int(uintptr(z.ZendValue.A))
	return *p
}

func (z *zVal) AsString() string {
	if z == nil {
		return ""
	}
	return C.GoStringN((*C.char)(z.ZendValue.A), C.int(uint32(uintptr(z.ZendValue.B))))
}

func (z *zVal) AsArray() api.ZArray {
	if z == nil {
		return nil
	}
	return NewZArray(z.ZendValue.A)
}

func (z *zVal) AsResource() api.ZResource {
	if z == nil {
		return nil
	}
	return NewZResource(z)
}

func (z *zVal) AsPointer() unsafe.Pointer {
	return nil
}

func (z *zVal) AsObject() api.ZObject {
	if z == nil {
		return nil
	}
	return NewZObject(z)
}

func NewZArray(p unsafe.Pointer) api.ZArray {
	if p == nil {
		return nil
	}
	return (*hashTable)(p)
}

func (z *hashTable) Len() int {
	if z == nil {
		return 0
	}
	return int(z.NumOfElements)
}

func (z *hashTable) ToSSMap() (m map[string]string) {
	m = make(map[string]string)
	forEach(z, func(keyNil bool, intKey int, stringKey string, zv api.ZVal) bool {
		if !keyNil && zv.IsString() {
			m[stringKey] = zv.AsString()
		}
		return true
	})
	return
}

func (z *hashTable) ToSZMap() (m map[string]api.ZVal) {
	m = make(map[string]api.ZVal)
	forEach(z, func(keyNil bool, intKey int, stringKey string, zv api.ZVal) bool {
		if !keyNil {
			m[stringKey] = zv
		}
		return true
	})
	return
}

func (z *hashTable) ToIIMap() (m map[interface{}]interface{}) {
	m = make(map[interface{}]interface{})
	forEach(z, func(keyNil bool, intKey int, stringKey string, zv api.ZVal) bool {
		if keyNil {
			m[intKey] = zv.Value()
		} else {
			m[stringKey] = zv.Value()
		}
		return true
	})
	return
}

func (z *hashTable) ToIZMap() (m map[interface{}]api.ZVal) {
	m = make(map[interface{}]api.ZVal)
	forEach(z, func(keyNil bool, intKey int, stringKey string, zv api.ZVal) bool {
		if keyNil {
			m[intKey] = zv
		} else {
			m[stringKey] = zv
		}
		return true
	})
	return
}

func forEach(z *hashTable, fn func(keyNil bool, intKey int, stringKey string, zv api.ZVal) bool) {
	if z == nil {
		return
	}
	z.InternalPointer = z.ListHead
	p := z.InternalPointer
	for {
		var (
			keyNil    bool
			intKey    int
			stringKey string
			zv        api.ZVal
		)
		if p == nil {
			break
		}
		if p.KeyLength == 0 {
			keyNil = true
			intKey = int(p.H)
		} else {
			stringKey = C.GoStringN((*C.char)(p.ArKey), C.int(p.KeyLength-1)) // 这里要减掉末尾的 0 字符。
		}
		zv = NewZVal(p.Data)
		if !fn(keyNil, intKey, stringKey, zv) {
			break
		}
		p = p.ListNext
	}
	return
}

func NewZResource(val *zVal) api.ZResource {
	if val == nil {
		return nil
	}
	return &zendResource{val}
}

func (z *zendResource) ID() int32 {
	if z == nil {
		return 0
	}
	return int32(uintptr(z.val.ZendValue.A))
}

func (z *zendResource) Pointer() unsafe.Pointer {
	if z == nil {
		return nil
	}
	return unsafe.Pointer(z.val)
}

func (z *zendResource) String() string {
	if z == nil {
		return "(nil zendResource)"
	}
	return fmt.Sprintf("zendResource(%d)", z.ID())
}

func NewZObject(val *zVal) api.ZObject {
	if val == nil {
		return nil
	}
	return &zendObject{val}
}

func (z *zendObject) ID() int32 {
	if z == nil {
		return 0
	}
	return int32(uintptr(z.val.ZendValue.A))
}

func (z *zendObject) Pointer() unsafe.Pointer {
	if z == nil {
		return nil
	}
	return unsafe.Pointer(z.val)
}

func (z *zendObject) String() string {
	if z == nil {
		return "(nil zendObject)"
	}
	return fmt.Sprintf("zendObject(%d)", z.ID())
}

func (z *zendString) String() string {
	if z == nil {
		return ""
	}
	return C.GoStringN((*C.char)(z.Val), C.int(z.Len))
}
