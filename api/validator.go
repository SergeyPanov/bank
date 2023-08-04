package api

import (
	"github.com/SergeyPanov/bank/util"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return util.IsSupportedCUrrency(currency)
	}

	return false
}
