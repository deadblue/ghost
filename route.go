package ghost

import (
	"github.com/deadblue/ghost/internal/context"
	"github.com/deadblue/ghost/internal/rule"
)

type _Handler func(Context) (View, error)

type _TreeKey struct {
	name  string
	depth int
}

type _TreeNode struct {
	seg      *rule.Segment
	hdr      _Handler
	children map[_TreeKey]*_TreeNode
}

type _RouteTable struct {
	// Static table
	st map[string]_Handler
	// Trees
	trs map[string]*_TreeNode
}

// Mount mounts controller into a request path which is described by name.
func (t *_RouteTable) Mount(name string, h _Handler) (err error) {
	// Parse rule
	r := &rule.Rule{}
	if err = r.FromMethodName(name); err != nil {
		return
	}

	if r.IsStatic {
		t.st[r.StaticKey()] = h
	} else {
		// TODO
	}
	return
}

func (t *_RouteTable) Resolve(ctx *context.Impl) _Handler {
	// Search static table
	sk := ctx.Method() + " " + ctx.Path()
	if h, found := t.st[sk]; found {
		return h
	}
	// TODO:
	return nil
}
