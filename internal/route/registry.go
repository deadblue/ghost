package route

import (
	"fmt"
	"github.com/deadblue/ghost/internal/route/method"
	"github.com/deadblue/ghost/internal/route/parser"
)

type Registry[T any] struct {
	// Strict table
	st _StrictTable[T]
	// Path tree
	pt _PathTrees[T]
}

func (r *Registry[T]) Init() *Registry[T] {
	r.st = _StrictTable[T]{}
	r.pt = _PathTrees[T]{}
	return r
}

func (r *Registry[T]) Mount(name string, root string, target T) (rs string, err error) {
	// Parse method name to rule
	rule := parser.Rule{IsStrict: true, Depth: 0}
	// Set path root
	setPathRoot(&rule, root)
	// Parse rule
	if err = method.Parse(name, &rule); err != nil {
		return
	}

	if rule.IsStrict {
		r.st.Put(rule.Method, rule.Path, target)
		rs = rule.String()
	} else {
		varMap := make(map[int]string)
		node := r.pt.GetTree(rule.Method, rule.Ext, rule.Depth)
		for ok := rule.Pieces.GoFirst(); ok; ok = rule.Pieces.Forward() {
			index, piece := rule.Pieces.Get()
			if piece.IsVar {
				node = node.GetVarChild()
				varMap[index] = piece.Name
			} else {
				node = node.GetChild(piece.Name)
			}
		}
		if node.path != "" {
			err = fmt.Errorf("conflict path: \"%s\" <=> \"%s\"", node.path, rule.Path)
		} else {
			// Setup leaf node
			node.path = rule.Path
			node.target = target
			node.varMap = varMap
			rs = rule.String()
		}
	}
	return
}
