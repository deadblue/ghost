package route

type _TreeNode[T any] struct {
	// Parent node
	parent *_TreeNode[T]
	// Children nodes
	nameNodes map[string]*_TreeNode[T]
	// Variable child
	varNode *_TreeNode[T]

	// Route target, only available on leaf node
	target T
	// Variables map, only available on leaf node
	varMap map[int]string
	// Full path, only available on leaf node
	path string
}

func (n *_TreeNode[T]) FindChild(key string) (child *_TreeNode[T], found bool) {
	child, found = n.nameNodes[key]
	return
}

func (n *_TreeNode[T]) GetChild(key string) *_TreeNode[T] {
	if child, found := n.FindChild(key); found {
		return child
	}
	return n.addChild(key)
}

func (n *_TreeNode[T]) addChild(key string) *_TreeNode[T] {
	if n.nameNodes == nil {
		n.nameNodes = make(map[string]*_TreeNode[T])
	}
	child := &_TreeNode[T]{
		parent: n,
	}
	n.nameNodes[key] = child
	return child
}

func (n *_TreeNode[T]) GetVarChild() *_TreeNode[T] {
	if n.varNode == nil {
		n.varNode = &_TreeNode[T]{
			parent: n,
		}
	}
	return n.varNode
}

type _TreeKey struct {
	method string
	ext    string
	depth  int
}

type _RouteTrees[T any] map[_TreeKey]*_TreeNode[T]

func (t _RouteTrees[T]) FindTree(method, ext string, depth int) (node *_TreeNode[T], found bool) {
	key := _TreeKey{
		method: method,
		ext:    ext,
		depth:  depth,
	}
	node, found = t[key]
	return
}

// GetTree tries to find an exists tree, and it will create one if not found.
func (t _RouteTrees[T]) GetTree(method, ext string, depth int) *_TreeNode[T] {
	key := _TreeKey{
		method: method,
		ext:    ext,
		depth:  depth,
	}
	if node, found := t[key]; found {
		return node
	} else {
		node = &_TreeNode[T]{}
		t[key] = node
		return node
	}
}
