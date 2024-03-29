package route

type PathVariableReceiver interface {
	SetPathVar(name, value string)
}

func (r *Registry[T]) Resolve(method, path string, pvr PathVariableReceiver) (target T) {
	// First, search in strict table
	var found bool
	if target, found = r.st.Get(method, path); found {
		return
	}

	// Split request path
	pieces, depth, ext := splitRequestPath(path)
	// Get tree entry
	var node *_TreeNode[T]
	if node, found = r.pt.FindTree(method, ext, depth); !found {
		return
	}
	// Search path on tree
	for i := 0; i < depth; {
		// Find child
		piece, child := pieces[i], (*_TreeNode[T])(nil)
		if child, found = node.FindChild(piece); !found {
			child = node.varNode
		}
		if child != nil {
			node, i = child, i+1
		} else {
			// Backtrace
			for node.parent != nil &&
				(node.parent.varNode == node || node.parent.varNode == nil) {
				node, i = node.parent, i-1
			}
			// When back to the root, that means we can not find a route
			if node.parent == nil {
				node = nil
				break
			} else {
				node = node.parent.varNode
			}
		}
	}
	if node != nil {
		if node.varMap != nil {
			for index, name := range node.varMap {
				value := pieces[index]
				pvr.SetPathVar(name, value)
			}
		}
		target = node.target
	}
	return
}
