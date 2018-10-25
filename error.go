package appmetrica

import (
	"errors"
	"strconv"
)

const prefix = `appmetrica: `

var ErrNotImplemented = errors.New(prefix + "not implemented")

func NewError(code int, message string) error {
	return errors.New(prefix + "[" + strconv.Itoa(code) + "] " + message)
}
