package token

import (
	"fmt"
	"strconv"
)

type Position struct {
	Filename string
	Line     int
	Column   int
}

func (p *Position) IsValid() bool { return p.Line > 0 }

func (p Position) String() string {
	s := p.Filename
	if p.IsValid() {
		if s != "" {
			s += ":"
		}
		s += strconv.Itoa(p.Line)
		if p.Column != 0 {
			s += fmt.Sprintf(":%d", p.Column)
		}
	}
	if s == "" {
		s = "-"
	}
	return s
}
