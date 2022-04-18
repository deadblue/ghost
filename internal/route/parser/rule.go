package parser

import (
	"github.com/deadblue/ghost/internal/container"
	"strings"
)

type Rule struct {
	// Strict flag
	IsStrict bool

	// Request method
	Method string
	// Request path
	Path string

	// Request path in pieces
	Pieces container.List[*PathPiece]
	// Request path depth
	Depth int
	// Request extensions
	Ext string
}

func (r *Rule) Init() {
	r.Pieces = container.List[*PathPiece]{}
}

func (r *Rule) String() string {
	size := len(r.Method) + len(r.Path) + 1
	if r.Ext != "" {
		size += len(r.Ext) + 1
	}

	buf := strings.Builder{}
	buf.Grow(size)
	buf.WriteString(r.Method)
	buf.WriteRune(' ')
	buf.WriteString(r.Path)
	if r.Ext != "" {
		buf.WriteRune('.')
		buf.WriteString(r.Ext)
	}
	return buf.String()
}
