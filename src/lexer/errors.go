package lexer

import "github.com/kkumar-gcc/enumgen/src/token"

type Error struct {
	Pos token.Position
	Msg string
}

func (r *Error) Error() string {
	if r.Pos.IsValid() {
		return r.Pos.String() + ": " + r.Msg
	}
	return r.Msg
}

type ErrorList []*Error

func (r *ErrorList) Add(pos token.Position, msg string) {
	*r = append(*r, &Error{Pos: pos, Msg: msg})
}

func (r *ErrorList) Reset() {
	*r = (*r)[0:0]
}

func (r ErrorList) Len() int {
	return len(r)
}
