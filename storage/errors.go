package storage

import (
	"errors"
)

const (
	BadAction      uint = iota
	BadPointID     uint = iota
	BadShieldID    uint = iota
	BadEachFunc    uint = iota
	InternalError  uint = iota
	NotFoundPoint  uint = iota
	NotFoundShield uint = iota
	ShieldExists   uint = iota
	Success        uint = iota

	HookKeyBad      uint = iota
	HookFuncBad     uint = iota
	HookErrorStruct uint = iota
	HookErrorPoint  uint = iota
	HookErrorShield uint = iota
)

func iotaToError(code uint, prefixs ...string) error {

	errorLine := ""

	switch code {
	case Success:
		return nil
	case InternalError:
		errorLine = "Internal error"
	case NotFoundShield:
		errorLine = "Shield not found"
	case NotFoundPoint:
		errorLine = "Point not found"
	case BadShieldID:
		errorLine = "Shield not found. Bad Shield ID."
	case BadPointID:
		errorLine = "Point not found. Bad point ID."
	case ShieldExists:
		errorLine = "Shield Exists"
	case BadAction:
		errorLine = "Action not found"
	case HookErrorShield:
		errorLine = "Runtime hook shield error"
	case HookErrorPoint:
		errorLine = "Runtime hook point error"
	case BadEachFunc:
		errorLine = "'Each' func has wrong type"
	default:
		errorLine = " same error....."
	}

	prefix := ""
	if len(prefixs) > 0 {
		prefix = prefixs[0] + " "
	}

	postfix := ""
	if len(prefixs) > 1 {
		postfix = " " + prefixs[1]
	}

	return errors.New(prefix + errorLine + postfix)
}
