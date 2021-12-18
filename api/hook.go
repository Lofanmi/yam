package api

import (
	"context"
	"unsafe"
)

type (
	Hook interface {
		BeforeExecuteInternal(ctx context.Context, executeData unsafe.Pointer, returnValue unsafe.Pointer)
		AfterExecuteInternal(ctx context.Context, executeData unsafe.Pointer, returnValue unsafe.Pointer)

		BeforeExecuteEx(ctx context.Context, executeData unsafe.Pointer)
		AfterExecuteEx(ctx context.Context, executeData unsafe.Pointer)
	}

	HookCallback interface {
		Before(ctx context.Context, data ExecuteData, returnValue ZVal)
		After(ctx context.Context, data ExecuteData, returnValue ZVal)
	}
)
