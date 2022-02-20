package php7

import "C"
import (
	"fmt"
	"unsafe"

	"github.com/Lofanmi/yam/api"
)

var _ api.ZVal = (*zVal)(nil)

type ZType byte

const (
	IsUndef       ZType = 0 /* regular data types */
	IsNull        ZType = 1
	IsFalse       ZType = 2
	IsTrue        ZType = 3
	IsLong        ZType = 4
	IsDouble      ZType = 5
	IsString      ZType = 6
	IsArray       ZType = 7
	IsObject      ZType = 8
	IsResource    ZType = 9
	IsReference   ZType = 10
	IsConstant    ZType = 11 /* constant expressions */
	IsConstantAst ZType = 12
	_IsBool       ZType = 13 /* fake types */
	IsCallable    ZType = 14
	IsIterable    ZType = 19
	IsVoid        ZType = 18
	IsIndirect    ZType = 15 /* internal types */
	IsPtr         ZType = 17
	_IsError      ZType = 20

	HashFlagPersistent        = 1 << 0 // #define HASH_FLAG_PERSISTENT          (1<<0)
	HashFlagApplyProtection   = 1 << 1 // #define HASH_FLAG_APPLY_PROTECTION    (1<<1)
	HashFlagPacked            = 1 << 2 // #define HASH_FLAG_PACKED              (1<<2)
	HashFlagInitialized       = 1 << 3 // #define HASH_FLAG_INITIALIZED         (1<<3)
	HashFlagStaticKeys        = 1 << 4 // #define HASH_FLAG_STATIC_KEYS         (1<<4) /* long and interned strings */
	HashFlagHasEmptyInd       = 1 << 5 // #define HASH_FLAG_HAS_EMPTY_IND       (1<<5)
	HashFlagAllowCowViolation = 1 << 6 // #define HASH_FLAG_ALLOW_COW_VIOLATION (1<<6)
)

type (
	zVal struct {
		ZendValue unsafe.Pointer
		U1        uint32
		U2        uint32
		// union {
		// 	uint32_t var_flags;
		// 	uint32_t next;       /* hash collision chain */  -> zend_array 保存冲突拉链的下一个元素
		// 	uint32_t cache_slot; /* literal cache slot */
		// 	uint32_t lineno;     /* line number (for ast nodes) */
		// } u2;
	}
	zendRefCountedH struct {
		RefCount uint32
		U        uint32
	}
	zendResource struct {
		Gc     zendRefCountedH
		Handle int32
		Type   int32
		Ptr    unsafe.Pointer
	}
	zendString struct {
		Gc  zendRefCountedH
		H   uint32
		Len int
		Val [1]byte
	}
	zendClassEntry struct {
		Type byte
		Name *zendString
		// 还有一堆字段没有定义，不影响使用，但是用的时候要十分小心。
	}
	zendFunction struct {
		// Type          byte       /* MUST be the first element of this struct! */
		// QuickArgFlags uint32
		Common struct {
			Type            byte    /* never used */
			ArgFlags        [3]byte /* bitset of arg_info.pass_by_reference */
			FnFlags         uint32
			FunctionName    *zendString
			Scope           *zendClassEntry // zend_class_entry
			Prototype       unsafe.Pointer  // zend_function
			NumArgs         uint32
			RequiredNumArgs uint32
			ArgInfo         unsafe.Pointer // zend_arg_info  /* index -1 represents the return value info, if any */
			Attributes      unsafe.Pointer // HashTable
		}
		// OpArray
		// InternalFunction
	}
	zendExecuteData struct {
		OpLine           unsafe.Pointer   // zend_op              /* executed opline                */
		Call             *zendExecuteData // zend_execute_data    /* current call                   */
		ReturnValue      *zVal
		Func             *zendFunction    // zend_function        /* executed function              */
		This             zVal             /* this + call_info + num_args    */
		PrevExecuteData  *zendExecuteData // zend_execute_data
		SymbolTable      unsafe.Pointer   // zend_array
		RuntimeCache     unsafe.Pointer   /* cache op_array->run_time_cache */
		ExtraNamedParams unsafe.Pointer   // zend_array
	}
	zendBucket struct {
		ZVal zVal
		H    uint        /* hash value (or numeric index)   */
		Key  *zendString /* string key or NULL for numerics */
	}
	zendArray struct {
		GC              zendRefCountedH
		U               uint32
		TableMask       uint32
		ArData          *zendBucket
		NumUsed         uint32
		NumOfElements   uint32
		TableSize       uint32
		InternalPointer uint32
		NextFreeElement int
		Destructor      unsafe.Pointer
	}
)

// ZendInlineHashFunc zend_inline_hash_func
func ZendInlineHashFunc(s string) (hash uint) {
	n := len(s)
	hash = 5381
	for i := 0; i < n; i++ {
		hash = (hash << 5) + hash + uint(s[i])
	}
	hash = hash | 0x8000000000000000
	return
}

func NewZVal(p unsafe.Pointer) api.ZVal {
	if p == nil {
		return nil
	}
	from := (*zVal)(p)
	return &zVal{
		ZendValue: from.ZendValue,
		U1:        from.U1,
		U2:        from.U2,
	}
}

func (z *zVal) Type() int {
	if z == nil {
		return 0
	}
	return int(z.U1)
}

func (z *zVal) Value() interface{} {
	if z == nil {
		return nil
	}
	switch ZType(z.U1) {
	case IsUndef:
		fallthrough
	case IsNull:
		return nil
	case IsFalse:
		return false
	case IsTrue:
		return true
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
	case IsPtr:
		return z.AsPointer()
	default:
		fmt.Printf("unknown type: %d\n", ZType(z.U1))
	}
	return nil
}

func (z *zVal) IsBool() bool {
	return ZType(z.U1) == IsTrue || ZType(z.U1) == IsFalse
}

func (z *zVal) IsInt() bool {
	return ZType(z.U1) == IsLong
}

func (z *zVal) IsString() bool {
	return ZType(z.U1) == IsString
}

func (z *zVal) IsFloat() bool {
	return ZType(z.U1) == IsDouble
}

func (z *zVal) IsArray() bool {
	return ZType(z.U1) == IsArray
}

func (z *zVal) IsResource() bool {
	return ZType(z.U1) == IsResource
}

func (z *zVal) IsPointer() bool {
	return ZType(z.U1) == IsPtr
}

func (z *zVal) AsInt() int {
	if z == nil {
		return 0
	}
	return int(uintptr(z.ZendValue))
}

func (z *zVal) AsBool() bool {
	if z == nil {
		return false
	}
	switch ZType(z.U1) {
	case IsTrue:
		return true
	default:
		return false
	}
}

func (z *zVal) AsFloat() float64 {
	if z == nil {
		return 0.0
	}
	p := new(float64)
	i := (*int)(unsafe.Pointer(p))
	*i = int(uintptr(z.ZendValue))
	return *p
}

func (z *zVal) AsString() string {
	if z == nil {
		return ""
	}
	return (*zendString)(z.ZendValue).String()
}

func (z *zVal) AsArray() api.ZArray {
	if z == nil {
		return nil
	}
	return NewZArray(z.ZendValue)
}

func (z *zVal) AsResource() api.ZResource {
	if z == nil {
		return nil
	}
	return NewZResource(z.ZendValue)
}

func (z *zVal) AsPointer() unsafe.Pointer {
	if z == nil {
		return nil
	}
	return z.ZendValue
}

func (z *zendString) String() string {
	if z == nil {
		return ""
	}
	return C.GoString((*C.char)(unsafe.Pointer(&z.Val[0])))
}

func NewZArray(p unsafe.Pointer) api.ZArray {
	if p == nil {
		return nil
	}
	return (*zendArray)(p)
}

func (z *zendArray) Len() int {
	if z == nil {
		return 0
	}
	return int(z.NumUsed)
}

func (z *zendArray) ToSSMap() (m map[string]string) {
	m = make(map[string]string)
	forEach(z, func(keyNil bool, intKey int, stringKey string, zv api.ZVal) bool {
		if !keyNil && zv.IsString() {
			m[stringKey] = zv.AsString()
		}
		return true
	})
	return
}

func (z *zendArray) ToSZMap() (m map[string]api.ZVal) {
	m = make(map[string]api.ZVal)
	forEach(z, func(keyNil bool, intKey int, stringKey string, zv api.ZVal) bool {
		if !keyNil {
			m[stringKey] = zv
		}
		return true
	})
	return
}

func (z *zendArray) ToIIMap() (m map[interface{}]interface{}) {
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

func (z *zendArray) ToIZMap() (m map[interface{}]api.ZVal) {
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

func forEach(z *zendArray, fn func(keyNil bool, intKey int, stringKey string, zv api.ZVal) bool) {
	if z == nil {
		return
	}
	// ZEND_HASH_FOREACH(_ht, indirect)
	for i := 0; i < z.Len(); i++ {
		bucket := (*zendBucket)(unsafe.Pointer(uintptr(unsafe.Pointer(z.ArData)) + uintptr(i)*unsafe.Sizeof(zendBucket{})))
		intKey := int(bucket.H)
		stringKey := bucket.Key.String()
		var zv api.ZVal
		zv = NewZVal(unsafe.Pointer(&(bucket.ZVal)))
		if ZType(bucket.ZVal.Type()) == IsIndirect {
			zv = NewZVal(bucket.ZVal.AsPointer())
		}
		if ZType(zv.(*zVal).Type()) == IsUndef {
			continue
		}
		if !fn(bucket.Key == nil, intKey, stringKey, zv) {
			break
		}
	}
	// ZEND_HASH_FOREACH_END
	return
}

func NewZResource(p unsafe.Pointer) api.ZResource {
	if p == nil {
		return nil
	}
	return (*zendResource)(p)
}

func (z *zendResource) ID() int32 {
	if z == nil {
		return 0
	}
	return z.Handle
}

func (z *zendResource) Pointer() unsafe.Pointer {
	if z == nil {
		return nil
	}
	return z.Ptr
}

func (z *zendResource) String() string {
	if z == nil {
		return "(nil zendResource)"
	}
	return fmt.Sprintf("zendResource(%d)", z.ID())
}
