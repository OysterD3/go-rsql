package main

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type Order int

const (
	OrderDescending Order = -1
	OrderAscending  Order = 1
)

type Filter struct {
	Field    string
	Operator Operator
	Value    interface{}
	Group    []*Filter
}

type Sort struct {
	Field string
	Order Order
}

type RSQL struct {
	Filter         []Filter
	Sort           []Sort
	NextCursor     string
	PreviousCursor string
	Limit          int
	Offset         int
	filterMap      map[string]interface{}
}

func (r RSQL) ToGORMScope(db *gorm.DB) *gorm.DB {
	for _, sort := range r.Sort {
		if sort.Order == OrderDescending {
			db.Order(fmt.Sprintf("%s desc", sort.Field))
		} else {
			db.Order(sort.Field)
		}
	}
	var str strings.Builder
	args := make([]interface{}, 0)
	for _, filter := range r.Filter {
		if filter.Group != nil && len(filter.Group) > 0 {
			str.WriteString("(")
			for _, f := range filter.Group {
				if f.Operator == OpOr || f.Operator == OpAnd {
					str.WriteString(fmt.Sprintf("%s ", sqlOperatorMap[f.Operator]))
					continue
				}
				str.WriteString(
					fmt.Sprintf("%s %s ? ", f.Field, sqlOperatorMap[f.Operator]),
				)
				args = append(args, f.Value)
			}
			str.WriteString(") ")
			continue
		}
		if filter.Operator == OpOr || filter.Operator == OpAnd {
			str.WriteString(fmt.Sprintf("%s ", sqlOperatorMap[filter.Operator]))
			continue
		}
		str.WriteString(
			fmt.Sprintf("%s %s ? ", filter.Field, sqlOperatorMap[filter.Operator]),
		)
		args = append(args, filter.Value)
	}
	db.Where(str.String(), args...)
	return db
}
