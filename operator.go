package main

type Operator string

const (
	OpEqual        Operator = "EQUAL"
	OpNotEqual     Operator = "NOT_EQUAL"
	OpLessThan     Operator = "LESS_THAN"
	OpGreaterThan  Operator = "GREATER_THAN"
	OpLessEqual    Operator = "LESS_EQUAL"
	OpGreaterEqual Operator = "GREATER_EQUAL"
	OpIn           Operator = "IN"
	OpNotIn        Operator = "NOT_IN"
	OpLike         Operator = "LIKE"
	OpNotLike      Operator = "NOT_LIKE"

	OpAnd Operator = "AND"
	OpOr  Operator = "OR"
)

var comparisonOperatorMap = map[string]Operator{
	"==":    OpEqual,
	"=eq=":  OpEqual,
	"!=":    OpNotEqual,
	"=ne=":  OpNotEqual,
	"<":     OpLessThan,
	"=lt=":  OpLessThan,
	"<=":    OpLessEqual,
	"=lte=": OpLessEqual,
	"=le=":  OpLessEqual,
	">":     OpGreaterThan,
	"=gt=":  OpGreaterThan,
	">=":    OpGreaterEqual,
	"=gte=": OpGreaterEqual,
	"=ge=":  OpGreaterEqual,
	"=in=":  OpIn,
	"=nin=": OpNotIn,
}

var logicalOperatorMap = map[string]Operator{
	";":   OpAnd,
	"and": OpAnd,
	"or":  OpOr,
	",":   OpOr,
}

var sqlOperatorMap = map[Operator]string{
	OpEqual:        "=",
	OpNotEqual:     "!=",
	OpLessThan:     "<",
	OpLessEqual:    "<=",
	OpGreaterThan:  ">",
	OpGreaterEqual: ">=",
	OpIn:           "IN",
	OpNotIn:        "NOT IN",
	OpLike:         "LIKE",
	OpNotLike:      "NOT LIKE",
	OpAnd:          "AND",
	OpOr:           "OR",
}
