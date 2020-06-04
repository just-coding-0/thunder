package thunder

type nodeType uint8

const (
	static nodeType = iota // default
	root                   // 根节点
)

type methodTrees []methodTree

type methodTree struct {
	method string
	root   *node
}

func (trees methodTrees) get(method string) *node {
	for _, tree := range trees {
		if tree.method == method {
			return tree.root
		}
	}
	return nil
}

type node struct {
	path     string        // 当前路径
	indices  string        // 索引
	priority uint32        // 权重
	nType    nodeType      // 节点类型
	children []*node       // 子节点
	handlers HandlersChain // 处理链
	fullPath string        // 完整路径
}

// 增加子节点
func (n *node) addRoute(path string, handlers HandlersChain) {
	fullPath := path
	n.priority++

	// Empty tree
	if len(n.path) == 0 && len(n.children) == 0 {
		n.insertChild(path, fullPath, handlers)
		n.nType = root
		return
	}

	parentFullPathIndex := 0

walk:
	for {
		i := longestCommonPrefix(path, n.path)

		// 最短公共前缀
		if i < len(n.path) {
			child := node{
				path:     n.path[i:],
				indices:  n.indices,
				children: n.children,
				handlers: n.handlers,
				priority: n.priority - 1,
				fullPath: n.fullPath,
			}

			n.children = []*node{&child}
			n.indices = BytesToString([]byte{n.path[i]})
			n.path = path[:i]
			n.handlers = nil
			n.fullPath = fullPath[:parentFullPathIndex+i]
		}

		if i < len(path) {
			path = path[i:]
			c := path[0]

			for i, max := 0, len(n.indices); i < max; i++ {
				if c == n.indices[i] { // 如果已经有该公共前缀
					parentFullPathIndex += len(n.path)
					i = n.incrementChildPrio(i)
					n = n.children[i]
					continue walk
				}
			}

			n.indices += BytesToString([]byte{c}) // zero copy
			child := &node{
				fullPath: fullPath,
			}
			n.children = append(n.children, child)
			n.incrementChildPrio(len(n.indices) - 1)
			n = child

			n.insertChild(path, fullPath, handlers)
			return
		}

		if n.handlers != nil {
			panic("handlers are already registered for path '" + fullPath + "'")
		}
		n.handlers = handlers
		n.fullPath = fullPath
		return
	}
}

func (n *node) insertChild(path string, fullPath string, handlers HandlersChain) {
	n.path = path
	n.fullPath = fullPath
	n.handlers = handlers
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

// 最长公共前缀
func longestCommonPrefix(a, b string) int {
	i := 0
	max := min(len(a), len(b))
	for i < max && a[i] == b[i] {
		i++
	}
	return i
}

// Increments priority of the given child and reorders if necessary
func (n *node) incrementChildPrio(pos int) int {
	cs := n.children
	cs[pos].priority++
	prio := cs[pos].priority

	// Adjust position (move to front)
	newPos := pos
	for ; newPos > 0 && cs[newPos-1].priority < prio; newPos-- {
		// Swap node positions
		cs[newPos-1], cs[newPos] = cs[newPos], cs[newPos-1]

	}

	// Build new index char string
	if newPos != pos {
		n.indices = n.indices[:newPos] + // Unchanged prefix, might be empty
			n.indices[pos:pos+1] + // The index char we move
			n.indices[newPos:pos] + n.indices[pos+1:] // Rest without char at 'pos'
	}

	return newPos
}
