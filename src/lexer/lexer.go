package lexer

import (
	"unicode"
	"unicode/utf8"

	"github.com/kkumar-gcc/enumgen/src/token"
)

type Lexer struct {
	filename string
	src      []byte

	ch     rune
	pos    int
	rdPos  int
	line   int
	column int

	mode int
}

const (
	EOF         = -1
	CommentMode = 1 << iota
)

func New(filename string, src []byte, mode int) *Lexer {
	l := &Lexer{
		filename: filename,
		src:      src,
		line:     1,
		column:   0,
		ch:       ' ',
		mode:     mode,
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

	// set current byte position
	r.pos = r.rdPos

	// decode next rune
	ch, w := rune(r.src[r.rdPos]), 1
	if ch >= utf8.RuneSelf {
		var size int
		if ch, size = utf8.DecodeRune(r.src[r.rdPos:]); size > 0 {
			w = size
		}
	}
	r.rdPos += w
	r.ch = ch

	if r.ch == '\n' {
		r.line++
		r.column = 0
		return
	}
	if r.ch == '\t' {
		r.column += w
		return
	}
	r.column++
}

func (r *Lexer) skipWhitespace() {
	for r.ch == ' ' || r.ch == '\t' || r.ch == '\n' || r.ch == '\r' {
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
	pos = token.Position{Filename: r.filename, Line: r.line, Column: r.column}
	switch ch := r.ch; {
	case isLetter(ch):
		lit = r.lexIdentifier()
		tok = token.Lookup(lit)
		if len(lit) > 1 {
			tok = token.Lookup(lit)
		} else {
			tok = token.IDENT
		}
	case isDecimal(ch) || ch == '.' && isDecimal(rune(r.peek())):
		tok, lit = r.lexNumber()
	default:
		r.next()
		switch ch {
		case EOF:
			return pos, token.EOF, ""
		case '\n':
			r.next()
			goto scanAgain
		case '"':
			tok = token.STRING
			lit = r.lexString()
		case '\'':
			tok = token.CHAR
			lit = r.lexChar()
		case '/':
			if r.ch == '/' {
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
			tok = token.RBRACKET
		default:
			tok = token.ILLEGAL
			lit = string(ch)
		}
	}

	return
}

func (r *Lexer) lexIdentifier() string {
	pos := r.pos

	for rdPos, b := range r.src[r.rdPos:] {
		if 'a' <= b && b <= 'z' || 'A' <= b && b <= 'Z' || b == '_' || '0' <= b && b <= '9' {
			r.column++
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
