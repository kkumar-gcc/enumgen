package parser

import (
	"fmt"

	"github.com/kkumar-gcc/enumgen/src/ast"
	"github.com/kkumar-gcc/enumgen/src/lexer"
	"github.com/kkumar-gcc/enumgen/src/token"
)

type Parser struct {
	l   *lexer.Lexer
	err lexer.ErrorList

	pos token.Position
	tok token.Token
	lit string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.next()
	return p
}

func (p *Parser) next() {
	p.pos, p.tok, p.lit = p.l.Lex()
}

func (p *Parser) tokenIs(tok token.Token) bool {
	return p.tok == tok
}

func (p *Parser) expect(tok token.Token, msg string) bool {
	if p.tokenIs(tok) {
		p.next()
		return true
	}
	p.errorExpected(msg)
	return false
}

func (p *Parser) errorExpected(msg string) {
	p.err.Add(p.pos, fmt.Sprintf("expected %s at %v, got %s", msg, p.pos, p.tok))
}

func (p *Parser) Errors() lexer.ErrorList {
	return p.err
}

func (p *Parser) Parse() *ast.File {
	file := &ast.File{
		FileStart:    p.pos,
		Declarations: []ast.Decl{},
		Comments:     []*ast.CommentGroup{},
	}

	for !p.tokenIs(token.EOF) {
		if p.tokenIs(token.COMMENT) {
			cg := p.consumeComments()
			file.Comments = append(file.Comments, cg)
			continue
		}

		if p.tokenIs(token.ENUM) {
			decl := p.parseEnum()
			file.Declarations = append(file.Declarations, decl)

			// Skip any extra tokens until we're at a position to parse a new declaration
			for !p.tokenIs(token.EOF) && !p.tokenIs(token.ENUM) && !p.tokenIs(token.COMMENT) {
				p.next()
			}
		} else {
			p.errorExpected("enum declaration")
			p.next()
		}
	}

	file.FileEnd = p.pos
	return file
}

func (p *Parser) consumeComments() *ast.CommentGroup {
	group := &ast.CommentGroup{}
	for p.tokenIs(token.COMMENT) {
		group.Add(&ast.Comment{Slash: p.pos, Text: p.lit})
		p.next()
	}
	return group
}

// EnumDefinition ::= { Comment } 'enum' Identifier [ TypeSpec ] MemberList
func (p *Parser) parseEnum() *ast.EnumDefinition {
	enum := &ast.EnumDefinition{Doc: p.consumeComments()}
	if !p.expect(token.ENUM, "enum") {
		return enum
	}
	enum.EnumPos = p.pos

	if !p.tokenIs(token.IDENT) {
		p.errorExpected("identifier")
		return enum
	}
	enum.Name = ast.Ident{NamePos: p.pos, Name: p.lit}
	p.next()

	if !p.tokenIs(token.LBRACKET) {
		p.errorExpected("'['")
		return enum
	}
	enum.TypeSpec = p.parseTypeSpec()

	if !p.expect(token.COLON, "':' after enum declaration") {
		// Skip to next potential valid token
		for !p.tokenIs(token.EOF) && !p.tokenIs(token.SEMICOLON) && !p.tokenIs(token.ENUM) {
			p.next()
		}
		return enum
	}

	for {
		// Consume any comments before the member
		if p.tokenIs(token.COMMENT) {
			p.consumeComments()
			continue
		}

		if p.tokenIs(token.IDENT) {
			member := p.parseMember()

			switch p.tok {
			case token.COMMA:
				member.TermPos = p.pos
				enum.Members = append(enum.Members, member)
				p.next()
				continue

			case token.SEMICOLON:
				member.TermPos = p.pos
				enum.Members = append(enum.Members, member)
				p.next()
				return enum

			default:
				p.errorExpected("',' or ';' after enum member")
				enum.Members = append(enum.Members, member)
				// Try to recover by skipping to next semicolon or enum
				for !p.tokenIs(token.EOF) && !p.tokenIs(token.SEMICOLON) && !p.tokenIs(token.ENUM) {
					p.next()
				}
				if p.tokenIs(token.SEMICOLON) {
					p.next()
				}
				return enum
			}
		}

		// If we get here, we couldn't parse any more members
		if !p.tokenIs(token.SEMICOLON) && !p.tokenIs(token.EOF) {
			p.errorExpected("';' after enum declaration")
			// Try to recover by skipping to next semicolon or enum
			for !p.tokenIs(token.EOF) && !p.tokenIs(token.SEMICOLON) && !p.tokenIs(token.ENUM) {
				p.next()
			}
			if p.tokenIs(token.SEMICOLON) {
				p.next()
			}
		} else if p.tokenIs(token.SEMICOLON) {
			p.next()
		}

		return enum
	}
}

// TypeSpec ::= '[' Type { ',' Type } ']'
func (p *Parser) parseTypeSpec() *ast.TypeSpec {
	ts := &ast.TypeSpec{LbrackPos: p.pos}
	p.next()

	for {
		if !p.tokenIs(token.IDENT) {
			p.errorExpected("type identifier")
			break
		}
		tr := &ast.TypeRef{Name: ast.Ident{NamePos: p.pos, Name: p.lit}}
		p.next()

		for p.tokenIs(token.PERIOD) {
			p.next()
			if !p.tokenIs(token.IDENT) {
				p.errorExpected("identifier after dot")
				break
			}
			tr = &ast.TypeRef{
				Package: &tr.Name,
				DotPos:  p.pos,
				Name:    ast.Ident{NamePos: p.pos, Name: p.lit},
			}
			p.next()
		}

		ts.Types = append(ts.Types, tr)

		if !p.tokenIs(token.COMMA) {
			break
		}
		p.next()
	}

	if p.tokenIs(token.RBRACKET) {
		ts.RbrackPos = p.pos
		p.next()
	} else {
		p.errorExpected("]")
	}

	return ts
}

// MemberDefinition ::= { Comment } Identifier [ MemberAssignment ] [ Terminator ]
func (p *Parser) parseMember() *ast.MemberDefinition {
	lead := p.consumeComments()

	if !p.tokenIs(token.IDENT) {
		p.errorExpected("identifier")
		return nil
	}

	// Use saved position for member name
	m := &ast.MemberDefinition{
		Doc:  lead,
		Name: ast.Ident{NamePos: p.pos, Name: p.lit},
	}
	p.next()

	if p.tokenIs(token.ASSIGN) {
		m.AssignPos = p.pos
		p.next()

		if p.isLiteral(p.tok) {
			lit1 := &ast.BasicLit{ValuePos: p.pos, Kind: p.tok, Value: p.lit}
			p.next()

			if p.tokenIs(token.COLON) {
				colonPos := p.pos
				p.next()
				if !p.isLiteral(p.tok) {
					p.errorExpected("literal after ':'")
					return m
				}
				lit2 := &ast.BasicLit{ValuePos: p.pos, Kind: p.tok, Value: p.lit}
				m.Value = &ast.KeyValueExpr{Key: lit1, Colon: colonPos, Value: lit2}
				p.next()
			} else {
				m.Value = lit1
			}
		} else {
			p.errorExpected("literal")
		}
	}

	return m
}

func (p *Parser) isLiteral(tok token.Token) bool {
	switch tok {
	case token.INT, token.FLOAT, token.STRING, token.CHAR, token.IDENT:
		return true
	default:
		return false
	}
}
