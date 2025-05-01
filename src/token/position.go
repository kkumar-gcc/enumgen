package token

import "fmt"

type Position struct {
	Line   int // Line number in the source file
	Column int // Column number in the source file
}

func (p Position) String() string {
	if p.Line <= 0 && p.Column <= 0 {
		return "<unknown position>"
	}
	return fmt.Sprintf("%d:%d", p.Line, p.Column)
}

func (p *Position) IsValid() bool { return p.Line > 0 }
