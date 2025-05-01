package token

import (
	"fmt"
	"unicode"
)

type Token int

const (
	ILLEGAL Token = iota
	EOF
	COMMENT

	ASSIGN
	LBRACKET  // [
	COMMA     // ,
	PERIOD    // .
	RBRACKET  // ]
	SEMICOLON // ;
	COLON     // :
	SUB       // -

	literal_beg
	IDENT  // main
	INT    // 12345
	FLOAT  // 123.45
	IMAG   // 123.45i
	CHAR   // 'a'
	STRING // "abc"
	literal_end

	keyword_beg
	ENUM
	KIND
	IOTA
	VALUE
	TRUE
	FALSE
	keyword_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",

	EOF:     "EOF",
	COMMENT: "COMMENT",

	IDENT:  "IDENT",
	INT:    "INT",
	FLOAT:  "FLOAT",
	IMAG:   "IMAG",
	CHAR:   "CHAR",
	STRING: "STRING",

	LBRACKET:  "[",
	RBRACKET:  "]",
	COMMA:     ",",
	COLON:     ":",
	ASSIGN:    "=",
	PERIOD:    ".",
	SEMICOLON: ";",
	SUB:       "-",

	ENUM:  "enum",
	KIND:  "kind",
	IOTA:  "iota",
	VALUE: "value",
	TRUE:  "true",
	FALSE: "false",
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token, keyword_end-(keyword_beg+1))
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}

func (t Token) String() string {
	if t >= 0 && int(t) < len(tokens) {
		s := tokens[t]
		if s != "" {
			return s
		}
	}
	return fmt.Sprintf("Token(%d)", t)
}

func (t Token) IsLiteral() bool { return literal_beg < t && t < literal_end }

func (t Token) IsKeyword() bool { return keyword_beg < t && t < keyword_end }

func Lookup(ident string) Token {
	if tok, isKeyword := keywords[ident]; isKeyword {
		return tok
	}
	return IDENT
}

func IsKeyword(name string) bool {
	_, ok := keywords[name]
	return ok
}

func IsIdentifier(name string) bool {
	if name == "" || IsKeyword(name) {
		return false
	}
	for i, c := range name {
		if !unicode.IsLetter(c) && c != '_' && (i == 0 || !unicode.IsDigit(c)) {
			return false
		}
	}
	return true
}
