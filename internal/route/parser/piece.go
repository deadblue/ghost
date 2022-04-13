package parser

import "strings"

type PathPiece struct {
	IsVar bool
	Name  string
}

func (p *PathPiece) String() string {
	// Calculate name length
	length := len(p.Name)
	if p.IsVar {
		length += 2
	}
	// Early return
	if length == 0 {
		return ""
	}
	buf := &strings.Builder{}
	buf.Grow(length)
	if p.IsVar {
		buf.WriteRune('{')
	}
	_, _ = buf.WriteString(p.Name)
	if p.IsVar {
		buf.WriteRune('}')
	}
	return buf.String()
}
