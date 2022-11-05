package main

import (
	"fmt"
	"github.com/OysterD3/go-rsql/errcode"
)

func (lex *lexer) scanFilterToken() (*Token, error) {
	lex.skipNonTokens()
	var (
		tk  = new(Token)
		err error
	)

	if lex.r < 0 {
		tk.Type = EOF
		return tk, nil
	}

	switch lex.r {
	case '"', '\'':
		tk, err = lex.scanQuotedToken()
		break
	case '(', ')':
		tk, err = lex.scanParenthesisToken()
		break
	case '=', '!', '>', '<':
		tk, err = lex.scanComparisonOperatorToken()
		break
	default:
		if isAlphanumeric(lex.r) {
			tk, err = lex.scanAlphanumericToken()
			break
		}
		tk, err = lex.scanLogicalOperatorToken()
		break
	}

	if err != nil {
		return nil, err
	}

	if tk == nil {
		return nil, errcode.UnknownError
	}
	return tk, nil
}

func (lex *lexer) scanQuotedToken() (*Token, error) {
	startPosition := lex.position
	currentQuote := lex.r

	isEscaped := false
	lex.next()

	for lex.r > 0 && lex.r != currentQuote || (lex.r == currentQuote && isEscaped) {
		if lex.r == '\\' {
			isEscaped = true
		}
		lex.next()
	}

	lex.next()

	if lex.r < 0 {
		return nil, errcode.UnclosedQuote
	} else {
		tk := new(Token)
		tk.Type = VALUE
		tk.Value = string(lex.buffer[startPosition+1 : lex.position-1])

		return tk, nil
	}
}

func (lex *lexer) scanParenthesisToken() (*Token, error) {
	startPosition := lex.position
	lex.next()

	val := string(lex.buffer[startPosition:lex.position])

	tk := new(Token)
	tk.Type = OPEN_PARENTHESIS
	if val == ")" {
		tk.Type = CLOSED_PARENTHESIS
	}
	tk.Value = val

	return tk, nil
}

func (lex *lexer) scanComparisonOperatorToken() (*Token, error) {
	startPosition := lex.position
	op := ""
	for s := range comparisonOperatorMap {
		if string(lex.buffer[startPosition:startPosition+len(s)]) == s {
			op = s
			break
		}
	}
	if op == "" {
		return nil, errcode.UnknownOperator
	}
	lex.position += len(op)
	lex.r = rune(lex.buffer[lex.position])
	tk := new(Token)
	tk.Type = COMPARISON_OPERATOR
	tk.Value = comparisonOperatorMap[op]

	return tk, nil
}

func (lex *lexer) scanLogicalOperatorToken() (*Token, error) {
	startPosition := lex.position
	op := ""
	for s := range logicalOperatorMap {
		if string(lex.buffer[startPosition:startPosition+len(s)]) == s ||
			(len(lex.buffer) >= startPosition+len(s)+2 && string(lex.buffer[startPosition:startPosition+len(s)+2]) == fmt.Sprintf(" %s ", s)) {
			op = s
			break
		}
	}
	if op == "" {
		return nil, errcode.UnknownOperator
	}
	lex.position += len(op)
	lex.r = rune(lex.buffer[lex.position])
	tk := new(Token)
	tk.Type = LOGICAL_OPERATOR
	tk.Value = logicalOperatorMap[op]

	return tk, nil
}

func (lex *lexer) scanAlphanumericToken() (*Token, error) {
	startPosition := lex.position
	for isAlphanumeric(lex.r) {
		lex.next()
	}
	val := string(lex.buffer[startPosition:lex.position])
	tk := new(Token)
	tk.Value = val

	if lex.peekNextByte() == '=' {
		tk.Type = FIELD
	} else {
		tk.Type = VALUE
	}
	return tk, nil
}

func isAlphanumeric(r rune) bool {
	return 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || r == '_' || r == '$' || '0' <= r && r <= '9' || r == '-' || r == '*'
}

func (lex *lexer) scanSortTokens() ([]*Token, error) {
	tks := make([]*Token, 0)
	startPosition := lex.position

	for {
		if !isAlphanumeric(lex.r) || lex.r == ',' {
			tk := new(Token)
			tk.Type = SORT_FIELD
			tk.Value = string(lex.buffer[startPosition:lex.position])
			tks = append(tks, tk)

			if lex.r == ',' {
				startPosition = lex.position + 1
			}
		}
		if lex.r < 0 {
			break
		}
		lex.next()
	}

	return tks, nil
}
