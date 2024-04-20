package api

import (
	"github.com/MathHRM/simple_bank/util"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	currency, ok := fieldLevel.Field().Interface().(string)

	if ok {
		return util.IsSupportedCurrency(currency)
	}
	return false;
}