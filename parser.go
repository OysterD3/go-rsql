package main

import (
	"fmt"
	"github.com/OysterD3/go-rsql/errcode"
	"reflect"
	"strconv"
	"strings"
)

func MustParse(query []byte) *RSQL {
	rsql, err := parse(query)
	if err != nil {
		panic(err)
	}
	return rsql
}

func Parse(query []byte) (*RSQL, error) {
	return parse(query)
}

func parse(query []byte) (*RSQL, error) {
	split := strings.Split(string(query), "&")
	var (
		rsql = new(RSQL)
		err  error
	)
	for _, s := range split {
		keyword, v := getKeyword(s)
		if keyword == FilterKey {
			var filters []Filter
			filters, err = parseFilters([]byte(v))
			if err != nil {
				return nil, err
			}
			rsql.Filter = filters

		} else if keyword == SortKey {
			var sorts []Sort
			sorts, err = parseSorts([]byte(v))
			if err != nil {
				return nil, err
			}
			rsql.Sort = sorts
		} else if keyword == NextCursorKey {
			rsql.NextCursor = v
		} else if keyword == PreviousCursorKey {
			rsql.PreviousCursor = v
		} else if keyword == LimitKey {
			var limit int64
			limit, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, err
			}
			rsql.Limit = int(limit)
		} else if keyword == OffsetKey {
			var offset int64
			offset, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, err
			}
			rsql.Offset = int(offset)
		} else {
			return nil, errcode.UnknownKeyword
		}
	}
	return rsql, nil
}

func parseFilters(src []byte) ([]Filter, error) {
	var (
		lex     = New(src)
		filters = make([]Filter, 0)
		group   = make([]Filter, 0)
		tk      *Token
		err     error
		filter  = Filter{}
		isIn    = false
		isGroup = false
	)

	for tk, err = lex.scanFilterToken(); err == nil && tk.Type != EOF; tk, err = lex.scanFilterToken() {
		if isGroup && tk.Type != CLOSED_PARENTHESIS {
			parseGroup(tk, &group, &filter)
			continue
		}

		switch tk.Type {
		case COMPARISON_OPERATOR:
			filter.Operator = tk.Value.(Operator)
			if filter.Operator == OpIn {
				isIn = true
				continue
			}
			break
		case LOGICAL_OPERATOR:
			if !isIn {
				filter.Operator = tk.Value.(Operator)
				filters = append(filters, filter)
				filter = Filter{}
			}
			break
		case VALUE:
			if isIn && filter.Value == nil {
				filter.Value = make([]string, 0)
				filter.Value = append(filter.Value.([]string), tk.Value.(string))
			} else if isIn && reflect.TypeOf(filter.Value).Kind() == reflect.Slice {
				filter.Value = append(filter.Value.([]string), tk.Value.(string))
			} else if reflect.ValueOf(tk.Value).Kind() == reflect.String &&
				(strings.Contains(tk.Value.(string), "*") || strings.Contains(tk.Value.(string), "%")) {
				if filter.Operator == OpEqual {
					filter.Operator = OpLike
				} else if filter.Operator == OpNotEqual {
					filter.Operator = OpNotLike
				}
				filter.Value = tk.Value
			} else {
				filter.Value = tk.Value
			}
			break
		case OPEN_PARENTHESIS:
			if !isIn {
				isGroup = true
			}
			break
		case CLOSED_PARENTHESIS:
			if isGroup {
				filter.Group = make([]*Filter, len(group))
				for idx := range group {
					filter.Group[idx] = &group[idx]
				}
				filters = append(filters, filter)
				filter = Filter{}
				group = make([]Filter, 0)
			}
			isGroup = false
			isIn = false
			break
		case FIELD:
			filter.Field = tk.Value.(string)
			break
		}

		if isGroup {
			continue
		}

		if reflect.ValueOf(filter.Value).IsValid() && reflect.TypeOf(filter.Value).Kind() != reflect.Slice {
			filters = append(filters, filter)
			filter = Filter{}
			continue
		}
	}

	if err != nil {
		return nil, err
	}

	return filters, nil
}

func parseGroup(tk *Token, filters *[]Filter, filter *Filter) {
	isIn := false

	if filter != nil && filter.Operator == OpIn {
		isIn = true
	}

	switch tk.Type {
	case COMPARISON_OPERATOR:
		filter.Operator = tk.Value.(Operator)
		break
	case LOGICAL_OPERATOR:
		if !isIn {
			filter.Operator = tk.Value.(Operator)
			*filters = append(*filters, *filter)
			*filter = Filter{}
		}
		break
	case VALUE:
		if isIn && filter.Value == nil {
			filter.Value = make([]string, 0)
			filter.Value = append(filter.Value.([]string), tk.Value.(string))
		} else if isIn && reflect.TypeOf(filter.Value).Kind() == reflect.Slice {
			filter.Value = append(filter.Value.([]string), tk.Value.(string))
		} else if reflect.ValueOf(tk.Value).Kind() == reflect.String &&
			(strings.Contains(tk.Value.(string), "*") || strings.Contains(tk.Value.(string), "%")) {
			if filter.Operator == OpEqual {
				filter.Operator = OpLike
			} else if filter.Operator == OpNotEqual {
				filter.Operator = OpNotLike
			}
			filter.Value = tk.Value

		} else {
			filter.Value = tk.Value
		}
		break
	case FIELD:
		filter.Field = tk.Value.(string)
		break
	}
	if reflect.ValueOf(filter.Value).IsValid() && reflect.TypeOf(filter.Value).Kind() != reflect.Slice {
		*filters = append(*filters, *filter)
		*filter = Filter{}
	}

	return
}

func parseSorts(src []byte) ([]Sort, error) {
	sorts := make([]Sort, 0)
	lex := New(src)
	tks, err := lex.scanSortTokens()

	if err != nil {
		return nil, err
	}

	for _, tk := range tks {
		sort := Sort{}
		sort.Order = OrderAscending
		sort.Field = strings.TrimPrefix(tk.Value.(string), "-")
		if strings.HasPrefix(tk.Value.(string), "-") {
			sort.Order = OrderDescending
		}
		sorts = append(sorts, sort)
	}

	return sorts, nil
}

func getKeyword(s string) (Keyword, string) {
	key := Keyword(s[:strings.IndexRune(s, '=')])
	return key, formatKey(key)
}

func formatKey(key Keyword) string {
	return fmt.Sprintf("%s=", key)
}
