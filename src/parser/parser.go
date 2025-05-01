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

	curPos, peekPos token.Position
	curTok, peekTok token.Token
	curLit, peekLit string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.next()
	p.next()
	return p
}

func (p *Parser) next() {
	p.curPos, p.curTok, p.curLit = p.peekPos, p.peekTok, p.peekLit
	p.peekPos, p.peekTok, p.peekLit = p.l.Lex()
}

func (p *Parser) tokenIs(tok token.Token) bool {
	return p.curTok == tok
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
	p.err.Add(p.curPos, fmt.Sprintf("expected %s at %v, got %s", msg, p.curPos, p.curTok))
}

func (p *Parser) Errors() lexer.ErrorList {
	return p.err
}

func (p *Parser) Parse() *ast.File {
	file := &ast.File{
		FileStart:    p.curPos,
		Declarations: []ast.Decl{},
		Comments:     []*ast.CommentGroup{},
	}

	for !p.tokenIs(token.EOF) {
		if p.tokenIs(token.COMMENT) {
			cg := p.consumeComments()
			file.Comments = append(file.Comments, cg)
			continue
		}

		decl := p.parseEnum()
		file.Declarations = append(file.Declarations, decl)
	}

	file.FileEnd = p.curPos
	return file
}

func (p *Parser) consumeComments() *ast.CommentGroup {
	group := &ast.CommentGroup{}
	for p.tokenIs(token.COMMENT) {
		group.Add(&ast.Comment{Slash: p.curPos, Text: p.curLit})
		p.next()
	}
	return group
}

// EnumDefinition ::= { Comment } 'enum' Identifier [ TypeSpec ] MemberList
func (p *Parser) parseEnum() *ast.EnumDefinition {
	enum := &ast.EnumDefinition{
		Doc: p.consumeComments(),
	}

	if !p.expect(token.ENUM, "enum") {
		return enum
	}
	enum.EnumPos = p.curPos

	if !p.tokenIs(token.IDENT) {
		p.errorExpected("identifier")
		return enum
	}
	enum.Name = ast.Ident{NamePos: p.curPos, Name: p.curLit}
	p.next()

	if p.tokenIs(token.LBRACKET) {
		enum.TypeSpec = p.parseTypeSpec()
	}

	for p.tokenIs(token.IDENT) || p.tokenIs(token.COMMENT) {
		m := p.parseMember()
		enum.Members = append(enum.Members, m)
	}

	return enum
}

// TypeSpec ::= '[' Type { ',' Type } ']'
func (p *Parser) parseTypeSpec() *ast.TypeSpec {
	ts := &ast.TypeSpec{LbrackPos: p.curPos}
	p.next()

	for {
		if !p.tokenIs(token.IDENT) {
			p.errorExpected("type identifier")
			break
		}
		tr := &ast.TypeRef{Name: ast.Ident{NamePos: p.curPos, Name: p.curLit}}
		p.next()

		for p.tokenIs(token.PERIOD) {
			p.next()
			if !p.tokenIs(token.IDENT) {
				p.errorExpected("identifier after dot")
				break
			}
			tr = &ast.TypeRef{
				Package: &tr.Name,
				DotPos:  p.curPos,
				Name:    ast.Ident{NamePos: p.curPos, Name: p.curLit},
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
		ts.RbrackPos = p.curPos
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
	m := &ast.MemberDefinition{
		Doc:  lead,
		Name: ast.Ident{NamePos: p.curPos, Name: p.curLit},
	}
	p.next()

	// MemberAssignment ::= '=' ( Literal | KeyValue )
	if p.tokenIs(token.ASSIGN) {
		m.AssignPos = p.curPos
		p.next()

		if p.isLiteral(p.curTok) {
			lit1 := &ast.BasicLit{ValuePos: p.curPos, Kind: p.curTok, Value: p.curLit}
			p.next()

			if p.tokenIs(token.COLON) {
				colon := p.curPos
				p.next()
				if !p.isLiteral(p.curTok) {
					p.errorExpected("literal after ':'")
					return m
				}
				lit2 := &ast.BasicLit{ValuePos: p.curPos, Kind: p.curTok, Value: p.curLit}
				m.Value = &ast.KeyValueExpr{Key: lit1, Colon: colon, Value: lit2}
				p.next()
			} else {
				m.Value = lit1
			}
		} else {
			p.errorExpected("literal")
		}
	}

	// Terminator ::= ',' | ';'
	if p.tokenIs(token.COMMA) || p.tokenIs(token.SEMICOLON) {
		m.TermPos = p.curPos
		p.next()
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
