package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// task business-level http error codes.
// the taskNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	taskNO       = 71
	taskName     = "task"
	taskBaseCode = errcode.HCode(taskNO)

	ErrCreateTask     = errcode.NewError(taskBaseCode+1, "failed to create "+taskName)
	ErrDeleteByIDTask = errcode.NewError(taskBaseCode+2, "failed to delete "+taskName)
	ErrUpdateByIDTask = errcode.NewError(taskBaseCode+3, "failed to update "+taskName)
	ErrGetByIDTask    = errcode.NewError(taskBaseCode+4, "failed to get "+taskName+" details")
	ErrListTask       = errcode.NewError(taskBaseCode+5, "failed to list of "+taskName)

	// error codes are globally unique, adding 1 to the previous error code
)
