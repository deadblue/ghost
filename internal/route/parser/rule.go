package parser

import (
	"github.com/deadblue/ghost/internal/container"
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
	return r.Method + " " + r.Path
}
