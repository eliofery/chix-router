package main

import (
	"errors"
	"github.com/eliofery/go-chix"
)

// Example пример реализации middleware
func Example() chix.Handler {
	return func(ctx *chix.Ctx) error {
		if false {
			return errors.New("некая ошибка")
		}

		return ctx.Next()
	}
}
