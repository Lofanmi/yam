package php7

//
// /*
// #include <stdlib.h>
// #include "extension.h"
// */
// import "C"
// import (
// 	"context"
// 	"log"
// 	"time"
// 	"unsafe"
//
// 	"github.com/Lofanmi/yam/api"
// 	"github.com/Lofanmi/yam/service"
//
// 	"github.com/opentracing/opentracing-go"
// )
//
// type (
// 	timeKey struct{}
// 	spanKey struct{}
// )
//
// //export ModuleStartup
// func ModuleStartup(pType int, moduleNumber int) int {
// 	service.InitServices()
// 	api.Globals = globals
// 	api.GlobalByKey = globalByKey
// 	api.CurlGetInfo = curlGetInfo
// 	return 0
// }
//
// //export ModuleShutdown
// func ModuleShutdown(pType int, moduleNumber int) int {
// 	service.CloseServices()
// 	return 0
// }
//
// //export RequestStartup
// func RequestStartup(pType int, moduleNumber int) int {
// 	span := service.Tracer().StartSpan("root")
// 	ctx := opentracing.ContextWithSpan(context.Background(), span)
// 	ctx = context.WithValue(ctx, spanKey{}, span)
// 	ctx = context.WithValue(ctx, timeKey{}, time.Now())
// 	service.ContextManager().SetContext(ctx)
// 	return 0
// }
//
// //export RequestShutdown
// func RequestShutdown(pType int, moduleNumber int) int {
// 	ctx := service.ContextManager().GetContext()
// 	t := ctx.Value(timeKey{}).(time.Time)
// 	since := time.Since(t)
// 	log.Printf("[yam] 请求结束，耗时：%s", since)
// 	span := ctx.Value(spanKey{}).(opentracing.Span)
// 	span.Finish()
// 	return 0
// }
//
// //export get_module
// func get_module() unsafe.Pointer {
// 	name := service.ExtensionName()
// 	version := service.ExtensionVersion()
// 	cName := C.CString(name)
// 	cVersion := C.CString(version)
// 	defer func() {
// 		C.free(unsafe.Pointer(cName))
// 		C.free(unsafe.Pointer(cVersion))
// 	}()
// 	m := C.yam_get_module(cName, cVersion)
// 	return unsafe.Pointer(m)
// }
//
// //export BeforeExecuteEx
// func BeforeExecuteEx(executeData unsafe.Pointer) {
// 	service.Hook().BeforeExecuteEx(service.ContextManager().GetContext(), executeData)
// }
//
// //export AfterExecuteEx
// func AfterExecuteEx(executeData unsafe.Pointer) {
// 	service.Hook().AfterExecuteEx(service.ContextManager().GetContext(), executeData)
// }
//
// //export BeforeExecuteInternal
// func BeforeExecuteInternal(executeData unsafe.Pointer, returnValue unsafe.Pointer) {
// 	service.Hook().BeforeExecuteInternal(service.ContextManager().GetContext(), executeData, returnValue)
// }
//
// //export AfterExecuteInternal
// func AfterExecuteInternal(executeData unsafe.Pointer, returnValue unsafe.Pointer) {
// 	service.Hook().AfterExecuteInternal(service.ContextManager().GetContext(), executeData, returnValue)
// }
//
// func globals(typ api.GlobalTrackVars, callback func(ret unsafe.Pointer)) {
// 	var s string
// 	switch typ {
// 	case api.GlobalTrackVarsPost:
// 		s = "_POST"
// 	case api.GlobalTrackVarsGet:
// 		s = "_GET"
// 	case api.GlobalTrackVarsCookie:
// 		s = "_COOKIE"
// 	case api.GlobalTrackVarsServer:
// 		s = "_SERVER"
// 	case api.GlobalTrackVarsEnv:
// 		s = "_ENV"
// 	case api.GlobalTrackVarsFiles:
// 		s = "_FILES"
// 	case api.GlobalTrackVarsRequest:
// 		s = "_REQUEST"
// 	default:
// 		return
// 	}
// 	sp := C.CString(s)
// 	defer func() {
// 		C.free(unsafe.Pointer(sp))
// 	}()
// 	C.yam_globals(sp, C.ulonglong(uint64(uintptr(unsafe.Pointer(&callback)))))
// }
//
// //export globalsCallback
// func globalsCallback(ret unsafe.Pointer, fn unsafe.Pointer) {
// 	callback := *((*func(pointer unsafe.Pointer))(fn))
// 	callback(ret)
// }
//
// func globalByKey(typ api.GlobalTrackVars, key string, callback func(ret unsafe.Pointer)) {
//
// }
//
// func curlGetInfo(ex unsafe.Pointer, callback func(ret unsafe.Pointer)) {
// 	// ch := unsafe.Pointer(uintptr(data.Pointer()) + 10*unsafe.Sizeof(0))
// 	C.yam_curl_getinfo(ex, C.ulonglong(uint64(uintptr(unsafe.Pointer(&callback)))))
// }
//
// //export curlGetInfoCallback
// func curlGetInfoCallback(ret unsafe.Pointer, fn unsafe.Pointer) {
// 	callback := *((*func(pointer unsafe.Pointer))(fn))
// 	callback(ret)
// }
