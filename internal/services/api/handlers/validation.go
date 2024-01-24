package handlers

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/rarime-link-svc/internal/data"
)

var (
	ErrOperator        = validation.NewError("validation_is_operator", "must be a valid operator")
	ValidationOperator = validation.NewStringRuleWithError(isOperator, ErrOperator)
)

func isOperator(op string) bool {
	_, ok := data.ProofOperatorFromString(op)
	return ok
}
