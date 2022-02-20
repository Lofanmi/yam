package sdks

import "C"
import (
	"unsafe"

	"github.com/Lofanmi/yam/api"
	"github.com/Lofanmi/yam/internal/hooks/curl"
	"github.com/Lofanmi/yam/internal/hooks/mysqli"
	"github.com/Lofanmi/yam/internal/hooks/redis"
	"github.com/Lofanmi/yam/internal/zend"
)

type hook struct {
	request   api.PHPRequest
	callbacks map[string]api.HookCallback
}

// NewHook 创建全局钩子，截获 PHP 的请求生命周期。
func NewHook(request api.PHPRequest) api.Hook {
	h := &hook{request: request}
	h.initCallbacks()
	return h
}

func (h *hook) initCallbacks() {
	h.callbacks = map[string]api.HookCallback{
		// cURL
		"curl_init()":         curl.NewInit(h.request),
		"curl_setopt()":       curl.NewSetOption(h.request),
		"curl_setopt_array()": curl.NewSetOptionArray(h.request),
		"curl_exec()":         curl.NewExec(h.request),
		"curl_close()":        curl.NewClose(h.request),
		// PDO
		// mysqli
		"mysqli::__construct()":  mysqli.NewObjectConstruct(h.request),
		"mysqli::mysqli()":       mysqli.NewObjectConstruct(h.request),
		"mysqli::query()":        mysqli.NewObjectQuery(h.request),
		"mysqli::prepare()":      mysqli.NewObjectPrepare(h.request),
		"mysqli_stmt::execute()": mysqli.NewObjectStmtExecute(h.request),
		// mysql_*
		"mysqli_connect()":      mysqli.NewConnect(h.request),
		"mysqli_query()":        mysqli.NewQuery(h.request),
		"mysqli_prepare()":      mysqli.NewPrepare(h.request),
		"mysqli_stmt_execute()": mysqli.NewStmtExecute(h.request),
		// Redis
		"Redis::*": redis.NewPHPRedis(h.request),
		// RedisCluster
		"RedisCluster::*": redis.NewPHPRedisCluster(h.request),
		// predis
		"Predis\\Client::*": redis.NewPRedis(h.request),
	}
}

func (h *hook) BeforeExecuteInternal(executeData unsafe.Pointer, pReturnValue unsafe.Pointer) {
	returnValue := zend.NewZVal(pReturnValue)
	data := zend.NewExecuteData(executeData)
	if callback, ok := h.callbacks[data.ClassFuncName()]; ok {
		callback.Before(data, returnValue)
		return
	}
	classFunctionAll := data.ClassName() + "::*"
	if callback, ok := h.callbacks[classFunctionAll]; ok {
		callback.Before(data, returnValue)
	}
}

func (h *hook) AfterExecuteInternal(executeData unsafe.Pointer, pReturnValue unsafe.Pointer) {
	returnValue := zend.NewZVal(pReturnValue)
	data := zend.NewExecuteData(executeData)
	if callback, ok := h.callbacks[data.ClassFuncName()]; ok {
		callback.After(data, returnValue)
		return
	}
	classFunctionAll := data.ClassName() + "::*"
	if callback, ok := h.callbacks[classFunctionAll]; ok {
		callback.After(data, returnValue)
	}
}

func (h *hook) BeforeExecuteEx(executeData unsafe.Pointer) {
	h.request.GetContext()
	data := zend.NewExecuteData(executeData)
	if callback, ok := h.callbacks[data.ClassFuncName()]; ok {
		callback.Before(data, nil)
		return
	}
	classFunctionAll := data.ClassName() + "::*"
	if callback, ok := h.callbacks[classFunctionAll]; ok {
		callback.Before(data, nil)
	}
}

func (h *hook) AfterExecuteEx(executeData unsafe.Pointer) {
	data := zend.NewExecuteData(executeData)
	if callback, ok := h.callbacks[data.ClassFuncName()]; ok {
		callback.After(data, nil)
		return
	}
	classFunctionAll := data.ClassName() + "::*"
	if callback, ok := h.callbacks[classFunctionAll]; ok {
		callback.After(data, nil)
	}
}
