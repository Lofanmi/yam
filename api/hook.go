package api

import (
	"unsafe"
)

type (
	Hook interface {
		BeforeExecuteInternal(executeData unsafe.Pointer, returnValue unsafe.Pointer)
		AfterExecuteInternal(executeData unsafe.Pointer, returnValue unsafe.Pointer)

		BeforeExecuteEx(executeData unsafe.Pointer)
		AfterExecuteEx(executeData unsafe.Pointer)
	}

	HookCallback interface {
		Before(data ExecuteData, returnValue ZVal)
		After(data ExecuteData, returnValue ZVal)
	}
)
