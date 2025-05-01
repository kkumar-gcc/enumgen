package ast

import "github.com/kkumar-gcc/enumgen/src/token"

type Node interface {
	Pos() token.Position
	End() token.Position
	String() string
}

type Expr interface {
	Node
	exprNode()
}

type Decl interface {
	Node
	declNode()
}

type Comment struct {
	Slash token.Position
	Text  string
}

func (r *Comment) Pos() token.Position { return r.Slash }
func (r *Comment) End() token.Position {
	return token.Position{
		Line:   r.Slash.Line,
		Column: r.Slash.Column + len(r.Text),
	}
}
func (r *Comment) String() string { return r.Text }

type CommentGroup struct {
	List []*Comment
}

func (r *CommentGroup) Pos() token.Position {
	if len(r.List) == 0 {
		return token.Position{}
	}
	return r.List[0].Pos()
}

func (r *CommentGroup) End() token.Position {
	if len(r.List) == 0 {
		return token.Position{}
	}
	return r.List[len(r.List)-1].End()
}

func (r *CommentGroup) Add(c *Comment) {
	r.List = append(r.List, c)
}

func (r *CommentGroup) String() string {
	var out string
	for _, c := range r.List {
		out += c.String() + "\n"
	}
	return out
}

type (
	Ident struct {
		NamePos token.Position
		Name    string
	}

	BasicLit struct {
		ValuePos token.Position
		Kind     token.Token
		Value    string
	}

	KeyValueExpr struct {
		Key   Expr
		Colon token.Position
		Value Expr
	}

	TypeRef struct {
		Package *Ident         // Optional: package.Type
		DotPos  token.Position // Dot position if Package != nil
		Name    Ident
	}

	TypeSpec struct {
		LbrackPos token.Position
		Doc       *CommentGroup
		Types     []*TypeRef
		Commas    []token.Position
		RbrackPos token.Position
	}

	MemberDefinition struct {
		Doc       *CommentGroup
		Name      Ident
		AssignPos token.Position
		Value     Expr
		TermPos   token.Position
	}

	// --- Declarations ---

	EnumDefinition struct {
		Doc      *CommentGroup
		EnumPos  token.Position
		Name     Ident
		TypeSpec *TypeSpec
		Members  []*MemberDefinition
	}

	BadDecl struct {
		From, To token.Position
	}

	File struct {
		Doc          *CommentGroup
		Declarations []Decl
		Comments     []*CommentGroup
		FileStart    token.Position
		FileEnd      token.Position
	}
)

func (r *Ident) Pos() token.Position { return r.NamePos }
func (r *Ident) End() token.Position {
	return token.Position{
		Line:   r.NamePos.Line,
		Column: r.NamePos.Column + len(r.Name),
	}
}
func (r *Ident) String() string { return r.Name }
func (r *Ident) exprNode()      {}

func (r *BasicLit) Pos() token.Position { return r.ValuePos }
func (r *BasicLit) End() token.Position {
	return token.Position{
		Line:   r.ValuePos.Line,
		Column: r.ValuePos.Column + len(r.Value),
	}
}
func (r *BasicLit) String() string { return r.Value }
func (r *BasicLit) exprNode()      {}

func (r *KeyValueExpr) Pos() token.Position { return r.Key.Pos() }
func (r *KeyValueExpr) End() token.Position { return r.Value.End() }
func (r *KeyValueExpr) String() string {
	return r.Key.String() + " : " + r.Value.String()
}
func (r *KeyValueExpr) exprNode() {}

func (r *TypeRef) Pos() token.Position {
	if r.Package != nil {
		return r.Package.Pos()
	}
	return r.Name.Pos()
}
func (r *TypeRef) End() token.Position { return r.Name.End() }
func (r *TypeRef) String() string {
	if r.Package != nil {
		return r.Package.String() + "." + r.Name.String()
	}
	return r.Name.String()
}
func (r *TypeRef) exprNode() {}

func (r *TypeSpec) Pos() token.Position { return r.LbrackPos }
func (r *TypeSpec) End() token.Position { return r.RbrackPos }
func (r *TypeSpec) String() string {
	out := "["
	for i, t := range r.Types {
		if i > 0 {
			out += ", "
		}
		out += t.String()
	}
	return out + "]"
}

func (r *MemberDefinition) Pos() token.Position { return r.Name.Pos() }
func (r *MemberDefinition) End() token.Position {
	if r.Value != nil {
		return r.Value.End()
	}
	return r.Name.End()
}
func (r *MemberDefinition) String() string {
	if r.Value != nil {
		return r.Name.String() + " = " + r.Value.String()
	}
	return r.Name.String()
}

func (r *EnumDefinition) Pos() token.Position { return r.EnumPos }
func (r *EnumDefinition) End() token.Position {
	if len(r.Members) > 0 {
		return r.Members[len(r.Members)-1].End()
	}
	return r.Name.End()
}
func (r *EnumDefinition) String() string {
	out := "enum " + r.Name.String()
	if r.TypeSpec != nil {
		out += " " + r.TypeSpec.String()
	}
	out += " {"
	for _, m := range r.Members {
		out += "\n  " + m.String()
	}
	out += "\n}"
	return out
}
func (r *EnumDefinition) declNode() {}

func (r *BadDecl) Pos() token.Position { return r.From }
func (r *BadDecl) End() token.Position { return r.To }
func (r *BadDecl) String() string {
	return "bad declaration from " + r.From.String() + " to " + r.To.String()
}
func (r *BadDecl) declNode() {}

func (r *File) Pos() token.Position { return r.FileStart }
func (r *File) End() token.Position {
	if len(r.Declarations) > 0 {
		return r.Declarations[len(r.Declarations)-1].End()
	}
	return r.FileEnd
}
func (r *File) String() string {
	var out string
	for _, decl := range r.Declarations {
		out += decl.String() + "\n"
	}
	return out
}
