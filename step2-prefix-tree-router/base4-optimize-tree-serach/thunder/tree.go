// Copyright 2020 just-codeding-0 . All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package thunder

import (
	"github.com/just-coding-0/thunder/internal/bytesconv"
	"net/url"
	"strings"
)

type nodeType uint8

const (
	static   nodeType = iota // default
	root                     // 根节点
	param                    // 参数节点
	catchAll                 // * 节点
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
	path      string        // 当前路径
	indices   string        // 索引
	priority  uint32        // 权重
	wildChild bool          // 是否为参数节点
	nType     nodeType      // 节点类型
	children  []*node       // 子节点
	handlers  HandlersChain // 处理链
	fullPath  string        // 完整路径
}

// 增加子节点
func (n *node) addRoute(path string, handlers HandlersChain) {
	fullPath := path
	n.priority++
	assert1(strings.Index(path, " ") == -1, "路径不能包含空格")

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
				path:      n.path[i:],
				indices:   n.indices,
				wildChild: n.wildChild,
				children:  n.children,
				handlers:  n.handlers,
				priority:  n.priority - 1,
				fullPath:  n.fullPath,
			}

			n.children = []*node{&child}
			n.indices = bytesconv.BytesToString([]byte{n.path[i]})
			n.path = path[:i]
			n.handlers = nil
			n.wildChild = false
			n.fullPath = fullPath[:parentFullPathIndex+i]
		}

		if i < len(path) {
			path = path[i:]

			if n.wildChild { // tip: 如果是通配符Node
				parentFullPathIndex += len(n.path)
				n = n.children[0]
				n.priority++

				// tip: 1.path大于当前节点,同时n.path == path[:len(n.path)]
				// tip: 2.*通配符节点与参数节点冲突
				// tip: 3.如果path和children.path相等(children不一定会有handler)或者还需要往后续匹配
				if len(path) >= len(n.path) && n.path == path[:len(n.path)] &&
					n.nType != catchAll &&
					(len(n.path) >= len(path) || path[len(n.path)] == '/') {
					continue walk
				}

				pathSeg := path
				if n.nType != catchAll {
					pathSeg = strings.SplitN(path, "/", 2)[0]
				}
				prefix := fullPath[:strings.Index(fullPath, pathSeg)] + n.path
				panic("'" + pathSeg +
					"' in new path '" + fullPath +
					"' conflicts with existing wildcard '" + n.path +
					"' in existing prefix '" + prefix +
					"'")

			}

			c := path[0]

			// 参数节点 first path
			if n.nType == param && c == '/' && len(n.children) == 1 {
				parentFullPathIndex += len(n.path)
				n = n.children[0]
				n.priority++ // n 增加权重
				continue walk
			}

			for i, max := 0, len(n.indices); i < max; i++ {
				if c == n.indices[i] { //tip: 如果已经有该公共前缀
					parentFullPathIndex += len(n.path)
					i = n.incrementChildPrio(i)
					n = n.children[i]
					continue walk
				}
			}

			// tip: 如果是普通节点直接插入
			if path[0] != '*' && path[0] != ':' {
				n.indices += bytesconv.BytesToString([]byte{c}) // zero copy
				child := &node{
					fullPath: fullPath,
				}
				n.children = append(n.children, child)
				// tip: 该函数会根据权重去排序索引
				n.incrementChildPrio(len(n.indices) - 1)
				n = child
			}

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
	for {
		// Find prefix until first wildcard
		wildcard, i, valid := findWildcard(path)
		if i < 0 { // No wildcard found
			break
		}

		// The wildcard name must not contain ':' and '*'
		if !valid {
			panic("only one wildcard per path segment is allowed, has: '" +
				wildcard + "' in path '" + fullPath + "'")
		}

		// check if the wildcard has a name
		if len(wildcard) < 2 {
			panic("wildcards must be named with a non-empty name in path '" + fullPath + "'")
		}

		// Check if this node has existing children which would be
		// unreachable if we insert the wildcard here
		if len(n.children) > 0 {
			panic("wildcard segment '" + wildcard +
				"' conflicts with existing children in path '" + fullPath + "'")
		}

		if wildcard[0] == ':' { // tip: 参数节点
			if i > 0 { // tip: wildcard前面还有路径
				n.path = path[:i]
				path = path[i:]
			}

			n.wildChild = true
			child := &node{
				nType:    param,
				path:     wildcard,
				fullPath: fullPath,
			}
			n.children = []*node{child}
			n = child
			n.priority++

			// tip: 如果path不是以/结尾,说明后面还有path
			if len(wildcard) < len(path) {
				path = path[len(wildcard):]

				child := &node{
					priority: 1,
					fullPath: fullPath,
				}
				n.children = []*node{child}
				n = child
				continue
			}

			// tip: 如果以/结尾,那么我们就插入了新节点
			n.handlers = handlers
			return
		}

		// catchAll
		if i+len(wildcard) != len(path) {
			panic("catch-all routes are only allowed at the end of the path in path '" + fullPath + "'")
		}

		if len(n.path) > 0 && n.path[len(n.path)-1] == '/' {
			panic("catch-all conflicts with existing handle for the path segment root in path '" + fullPath + "'")
		}

		// currently fixed width 1 for '/'
		i--
		if path[i] != '/' {
			panic("no / before catch-all in path '" + fullPath + "'")
		}

		n.path = path[:i]

		child := &node{
			wildChild: true,
			nType:     catchAll,
			fullPath:  fullPath,
		}

		n.children = []*node{child}
		n.indices = string('/')
		n = child
		n.priority++

		child = &node{
			path:     path[i:],
			nType:    catchAll,
			handlers: handlers,
			priority: 1,
			fullPath: fullPath,
		}
		n.children = []*node{child}

		return
	}

	n.path = path
	n.handlers = handlers
	n.fullPath = fullPath
}

type nodeValue struct {
	handlers HandlersChain
	params   *Params
	tsr      bool
	fullPath string
}

func (n *node) getValue(path string, params *Params, unescape bool) (value nodeValue) {
walk:
	for {
		prefix := n.path
		if len(path) > len(prefix) { // 如果path大于前缀
			if path[:len(prefix)] == prefix { // 前缀相等,直接寻找下一个子节点
				path = path[len(prefix):]
				if !n.wildChild { // tip: 如果当前节点不是通配符节点,则使用索引进行匹配
					idxc := path[0]
					for i, c := range []byte(n.indices) {
						if c == idxc {
							n = n.children[i]
							continue walk
						}
					}

					// not found 直接返回
					value.tsr = path == "/" && n.handlers != nil
					return
				}

				// 使用通配符子节点进行匹配
				n = n.children[0]
				switch n.nType {
				case param:
					// tip: 寻找最近的/通配符
					end := 0
					for end < len(path) && path[end] != '/' {
						end++
					}

					// Save param value
					if params != nil {
						if value.params == nil {
							value.params = params
						}
						// tip:在预分配slice中扩容
						i := len(*value.params)
						*value.params = (*value.params)[:i+1]
						val := path[:end]
						if unescape {
							if v, err := url.QueryUnescape(val); err == nil {
								val = v
							}
						}
						(*value.params)[i] = Param{
							Key:   n.path[1:],
							Value: val,
						}
					}

					// we need to go deeper!
					if end < len(path) {
						if len(n.children) > 0 {
							path = path[end:]
							n = n.children[0]
							continue walk
						}

						value.tsr = len(path) == end+1
						return
					}

					// 如果是以参数节点结束
					if value.handlers = n.handlers; value.handlers != nil {
						value.fullPath = n.fullPath
						return
					}
					if len(n.children) == 1 {
						// No handle found. Check if a handle for this path + a
						// trailing slash exists for TSR recommendation
						n = n.children[0]
						value.tsr = n.path == "/" && n.handlers != nil
					}
					return

				case catchAll:
					// Save param value
					if params != nil {
						if value.params == nil {
							value.params = params
						}
						// Expand slice within preallocated capacity
						// tip:在预分配slice中扩容
						i := len(*value.params)
						*value.params = (*value.params)[:i+1]
						val := path
						if unescape {
							if v, err := url.QueryUnescape(path); err == nil {
								val = v
							}
						}
						if len(val) > 0  && val[0] == '/'{
							val = val[1:]
						}

						(*value.params)[i] = Param{
							Key:   n.path[2:],
							Value: val,
						}
					}

					value.handlers = n.handlers
					value.fullPath = n.fullPath
					return

				default:
					panic("invalid node type")
				}
			}
		}

		if path == prefix {
			// We should have reached the node containing the handle.
			// Check if this node has a handle registered.
			// tip: 寻找到指定节点,返回
			if value.handlers = n.handlers; value.handlers != nil {
				value.fullPath = n.fullPath
				return
			}

			// If there is no handle for this route, but this route has a
			// wildcard child, there must be a handle for this path with an
			// additional trailing slash
			// tip:
			if path == "/" && n.wildChild && n.nType != root {
				value.tsr = true
				return
			}

			// No handle found. Check if a handle for this path + a
			// trailing slash exists for trailing slash recommendation
			for i, c := range []byte(n.indices) {
				if c == '/' {
					n = n.children[i]
					value.tsr = (len(n.path) == 1 && n.handlers != nil) ||
						(n.nType == catchAll && n.children[0].handlers != nil)
					return
				}
			}

			return
		}

		// Nothing found. We can recommend to redirect to the same URL with an
		// extra trailing slash if a leaf exists for that path
		value.tsr = (path == "/") ||
			(len(prefix) == len(path)+1 && prefix[len(path)] == '/' &&
				path == prefix[:len(prefix)-1] && n.handlers != nil)
		return
	}
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

func findWildcard(path string) (wildcard string, i int, valid bool) {
	// Find start
	for start, c := range []byte(path) {
		// A wildcard starts with ':' (param) or '*' (catch-all)
		if c != ':' && c != '*' {
			continue
		}

		// Find end and check for invalid characters
		valid = true
		for end, c := range []byte(path[start+1:]) {
			switch c {
			case '/':
				return path[start : start+1+end], start, valid
			case ':', '*':
				valid = false
			}
		}
		return path[start:], start, valid
	}
	return "", -1, false
}
