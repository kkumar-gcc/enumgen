package lexer

import (
	"unicode"
	"unicode/utf8"

	"github.com/kkumar-gcc/enumgen/src/token"
)

type Lexer struct {
	src []byte

	ch    rune
	pos   int
	rdPos int
	line  int

	insertSemi bool
	mode       int
}

const (
	EOF         = -1
	CommentMode = 1 << iota
)

func New(src []byte, mode int) *Lexer {
	l := &Lexer{
		src:  src,
		line: 1,
		ch:   ' ',
		mode: mode,
	}
	l.next()
	return l
}

func (r *Lexer) next() {
	if r.rdPos >= len(r.src) {
		r.pos = len(r.src)
		r.ch = EOF
		return
	}

	r.pos = r.rdPos
	s, w := rune(r.src[r.rdPos]), 1
	if s == 0 {
		// error
	}

	// Increment line number when newline is found
	if s == '\n' {
		r.line++
		//r.pos = 0
	}

	r.rdPos += w
	r.ch = s
}

func (r *Lexer) skipWhitespace() {
	for r.ch == ' ' || r.ch == '\t' || r.ch == '\n' && !r.insertSemi || r.ch == '\r' {
		r.next()
	}
}

func (r *Lexer) peek() byte {
	if r.rdPos < len(r.src) {
		return r.src[r.rdPos]
	}

	return 0
}

func (r *Lexer) Lex() (pos token.Position, tok token.Token, lit string) {
scanAgain:
	r.skipWhitespace()
	pos = token.Position{Line: r.line, Column: r.pos}
	insertSemi := false
	switch ch := r.ch; {
	case isLetter(ch):
		lit = r.lexIdentifier()
		tok = token.Lookup(lit)
		if len(lit) > 1 {
			tok = token.Lookup(lit)
		} else {
			tok = token.IDENT
			insertSemi = true
		}
	case isDecimal(ch) || ch == '.' && isDecimal(rune(r.peek())):
		insertSemi = true
		tok, lit = r.lexNumber()
	default:
		r.next()
		switch ch {
		case EOF:
			if r.insertSemi {
				r.insertSemi = false
				return pos, token.SEMICOLON, ";"
			}

			return pos, token.EOF, ""
		case '\n':
			r.insertSemi = false
			return pos, token.SEMICOLON, "\n"
		case '"':
			insertSemi = true
			tok = token.STRING
			lit = r.lexString()
		case '\'':
			insertSemi = true
			tok = token.CHAR
			lit = r.lexChar()
		case '/':
			if r.ch == '/' {
				insertSemi = r.insertSemi
				comment := r.lexComment()
				if r.mode&CommentMode == 0 {
					goto scanAgain
				}
				tok = token.COMMENT
				lit = comment
			} else {
				tok = token.ILLEGAL
				lit = "/"
			}
		case ',':
			tok = token.COMMA
		case ':':
			tok = token.COLON
		case '=':
			tok = token.ASSIGN
		case '-':
			tok = token.SUB
		case ';':
			tok = token.SEMICOLON
			lit = ";"
		case '[':
			tok = token.LBRACKET
		case ']':
			insertSemi = true
			tok = token.RBRACKET
		default:
			insertSemi = r.insertSemi
			tok = token.ILLEGAL
			lit = string(ch)
		}
	}

	r.insertSemi = insertSemi
	return
}

func (r *Lexer) lexIdentifier() string {
	pos := r.pos

	for rdPos, b := range r.src[r.rdPos:] {
		if 'a' <= b && b <= 'z' || 'A' <= b && b <= 'Z' || b == '_' || '0' <= b && b <= '9' {
			continue
		}

		r.rdPos += rdPos
		r.next()
		for isLetter(r.ch) || isDigit(r.ch) {
			r.next()
		}
		goto exit
	}

	r.pos = len(r.src)
	r.rdPos = len(r.src)
	r.ch = EOF
exit:
	return string(r.src[pos:r.pos])
}

func (r *Lexer) lexString() string {
	pos := r.pos - 1
	for {
		ch := r.ch
		if ch == '\n' || ch < 0 {
			// error
			break
		}

		r.next()
		if ch == '\\' {
			// Handle escape sequences
			if r.ch == '"' || r.ch == '\\' || r.ch == 'n' || r.ch == 't' || r.ch == 'r' {
				r.next()
				continue
			}
		} else if ch == '"' {
			break
		}
	}

	return string(r.src[pos:r.pos])
}

func (r *Lexer) lexChar() string {
	pos := r.pos - 1
	for {
		ch := r.ch
		if ch == '\n' || ch < 0 {
			// error
			break
		}

		r.next()
		if ch == '\\' {
			// Handle escape sequences
			if r.ch == '\'' || r.ch == '\\' || r.ch == 'n' || r.ch == 't' || r.ch == 'r' {
				r.next()
				continue
			}
		} else if ch == '\'' {
			break
		}
	}

	return string(r.src[pos:r.pos])
}

func (r *Lexer) lexNumber() (token.Token, string) {
	start := r.pos
	tok := token.ILLEGAL

	if r.ch == '.' {
		tok = token.FLOAT
		r.next()
		for isDigit(r.ch) {
			r.next()
		}
		return tok, string(r.src[start:r.pos])
	}

	tok = token.INT
	for isDigit(r.ch) {
		r.next()
	}

	if r.ch == '.' {
		tok = token.FLOAT
		r.next()
		for isDigit(r.ch) {
			r.next()
		}
	}

	return tok, string(r.src[start:r.pos])
}

func (r *Lexer) lexComment() string {
	pos := r.pos - 1
	if r.ch == '/' {
		r.next()
		for r.ch != '\n' && r.ch != EOF && r.ch >= 0 {
			r.next()
		}
	}

	// Find the position of the last non-whitespace character
	endPos := r.pos
	for i := r.pos - 1; i >= pos; i-- {
		if i < len(r.src) && (r.src[i] == ' ' || r.src[i] == '\t' || r.src[i] == '\r') {
			endPos = i
		} else {
			break
		}
	}

	return string(r.src[pos:endPos])
}

func isLetter(ch rune) bool {
	return 'a' <= lower(ch) && lower(ch) <= 'z' || ch == '_' || ch >= utf8.RuneSelf && unicode.IsLetter(ch)
}

func isDigit(ch rune) bool {
	return isDecimal(ch) || ch >= utf8.RuneSelf && unicode.IsDigit(ch)
}

func lower(ch rune) rune { return ('a' - 'A') | ch }

func isDecimal(ch rune) bool { return '0' <= ch && ch <= '9' }
func isHex(ch rune) bool     { return '0' <= ch && ch <= '9' || 'a' <= lower(ch) && lower(ch) <= 'f' }

func digitVal(ch rune) int {
	switch {
	case '0' <= ch && ch <= '9':
		return int(ch - '0')
	case 'a' <= lower(ch) && lower(ch) <= 'f':
		return int(lower(ch) - 'a' + 10)
	}
	return 16
}
