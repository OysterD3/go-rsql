package errcode

import "errors"

var (
	UnclosedQuote       = errors.New("closing quote is not found")
	UnclosedParenthesis = errors.New("closing parenthesis is not found")
	UnknownOperator     = errors.New("unknown operator")
	UnknownError        = errors.New("unknown error")
	UnknownIdentifier   = errors.New("unknown identifier")
	UnknownKeyword      = errors.New("unknown keyword")
)
