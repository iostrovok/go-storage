package storage

import (
	"errors"
)

const (
	BadPointID     uint = iota
	BadShieldID    uint = iota
	InternalError  uint = iota
	NotFoundPoint  uint = iota
	NotFoundShield uint = iota
	ShieldExists   uint = iota
	Success        uint = iota
	BadAction      uint = iota
)

func iotaToError(code uint, prefixs ...string) error {

	prefix := ""
	if len(prefixs) > 0 {
		prefix = prefixs[0] + ". "
	}

	switch code {
	case Success:
		return nil
	case InternalError:
		return errors.New(prefix + "Internal error")
	case NotFoundShield:
		return errors.New(prefix + "Shield not found")
	case NotFoundPoint:
		return errors.New(prefix + "Point not found")
	case BadShieldID:
		return errors.New(prefix + "Shield not found. Bad Shield ID.")
	case BadPointID:
		return errors.New(prefix + "Point not found. Bad point ID.")
	case ShieldExists:
		return errors.New(prefix + "Shield Exists")
	case BadAction:
		return errors.New(prefix + "Action not found")
	}
	return errors.New(prefix + " same error.....")
}
