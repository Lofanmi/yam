package php7

import (
	"unsafe"
)

type redisObject struct {
	Sock *redisSock
}

type redisSock struct {
	stream    unsafe.Pointer
	streamCtx unsafe.Pointer
	Host      *zendString
	Port      int32
}

type pdoDbhObject struct {
	inner unsafe.Pointer
}

type pdoDbh struct {
	methods            unsafe.Pointer // pdo_dbh_methods *methods
	driverData         unsafe.Pointer // driver_data
	username, password unsafe.Pointer
	reservedFlags      int32
	dataSource         unsafe.Pointer
	dataSourceLen      int
}
