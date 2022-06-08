package engine

type node struct {
	path     string
	fullPath string

	children []*node
	handlers HandlersChain
}

func (n *node) addRoute(path string, handlers HandlersChain) {
	fullPath := path
	if len(n.path) == 0 && len(n.children) == 0 {
		n.insertChild(path, fullPath, handlers)
	}

	parentFullPathIndex := 0

	for {
		pnum := longestPrefix(path, n.path) // pnum means the number of the longest prefix

		if pnum < len(n.path) {
			child := &node{
				path:     n.path[pnum:],
				fullPath: n.fullPath,
				children: n.children,
				handlers: n.handlers,
			}
			n.children = []*node{child}
			n.path = path[:pnum]
			n.fullPath = fullPath[:parentFullPathIndex+pnum]
		}

		if pnum < len(path) {
			path = path[pnum:]
		}

		n.insertChild(path, fullPath, handlers)
		return
	}
}

func (n *node) insertChild(path string, fullPath string, handlers HandlersChain) {
	child := &node{
		path:     path,
		fullPath: fullPath,
		children: nil,
		handlers: handlers,
	}
	n.addChild(child)
}

func (n *node) addChild(child *node) {
	n.children = append(n.children, child)
}

func (n *node) getNodeValue(fullPath string) *node {
	cnode := n
	path := fullPath
	c := path[0]
walk:
	for {
		if len(path) > len(cnode.path) {
			if path[0:len(cnode.path)] != cnode.path {
				return nil
			}
			path = path[len(cnode.path):]
			c = path[0]

			for _, child := range cnode.children {
				if c == child.path[0] {
					cnode = child
					break walk
				}
			}
			return nil
		}

		if path == cnode.path {
			return cnode
		}
		break
	}
	return nil
}

type methodTree struct {
	method string
	root   *node
}

type methodTrees []*methodTree

func (trees methodTrees) get(method string) *node {
	for _, tree := range trees {
		if tree.method == method {
			return tree.root
		}
	}
	return nil
}
