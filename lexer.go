package main

import "fmt"

// TokenType is a type for describing tokens mnemonically.
type TokenType int

type OperatorName string

type identifierType int

// Token represents a single token in the input stream.
// Type: mnemonic name (numeric).
// Value: string value of the token from the original stream.
type Token struct {
	Type     TokenType
	Value    interface{}
	Operator OperatorName
}

// Values for TokenType
const (
	EOF TokenType = iota
	IDENTIFIER
	FIELD
	COMPARISON_OPERATOR
	LOGICAL_OPERATOR
	VALUE
	SORT_FIELD
	OPEN_PARENTHESIS
	CLOSED_PARENTHESIS
)

var tokenTypes = [...]string{
	EOF:                 "EOF",
	IDENTIFIER:          "IDENTIFIER",
	FIELD:               "FIELD",
	COMPARISON_OPERATOR: "COMPARISON_OPERATOR",
	LOGICAL_OPERATOR:    "LOGICAL_OPERATOR",
	VALUE:               "VALUE",
	OPEN_PARENTHESIS:    "OPEN_PARENTHESIS",
	CLOSED_PARENTHESIS:  "CLOSED_PARENTHESIS",
}

func (tk Token) String() string {
	return fmt.Sprintf("Token{%s, '%s', %s}", tokenTypes[tk.Type], tk.Value, tk.Operator)
}

const (
	IdentifierFilter identifierType = iota
	IdentifierSort
)

type lexer struct {
	r        rune
	buffer   []byte
	position int
}

func New(i []byte) *lexer {
	l := new(lexer)
	l.position = 0
	l.buffer = i
	l.r = []rune(string(i))[0]

	return l
}

// next advances the lexer's internal state to point to the next rune in the
// input.
func (lex *lexer) next() {
	if lex.position < len(lex.buffer)-1 {
		lex.position += 1
		lex.r = rune(lex.buffer[lex.position])
	} else {
		lex.position = len(lex.buffer)
		lex.r = -1 // EOF
	}
}

// peekNextByte returns the next byte in the stream (the one after lex.r).
// Note: a single byte is peeked at - if there's a rune longer than a byte
// there, only its first byte is returned.
func (lex *lexer) peekNextByte() rune {
	if lex.position < len(lex.buffer)-1 {
		return rune(lex.buffer[lex.position])
	} else {
		return -1
	}
}

func (lex *lexer) scanAllFilterTokens() ([]*Token, error) {
	var (
		err error
		tks = make([]*Token, 0)
		tk  = new(Token)
	)
	for tk, err = lex.scanFilterToken(); err == nil && tk.Type != EOF; tk, err = lex.scanFilterToken() {
		tks = append(tks, tk)
	}

	return tks, err
}
