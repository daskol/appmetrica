package appmetrica

import (
	"errors"
	"strconv"
)

const prefix = `appmetrica: `

func NewError(code int, message string) error {
	return errors.New(prefix + "[" + strconv.Itoa(code) + "] " + message)
}
