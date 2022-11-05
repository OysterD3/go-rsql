package main

func (lex *lexer) skipNonTokens() {
	for lex.r == ' ' ||
		lex.r == '\t' ||
		lex.r == '\n' ||
		lex.r == '\r' {
		lex.next()
	}
}
