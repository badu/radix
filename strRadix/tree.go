package strRadix

import "strings"

func (t *Tree) compare(what string) bool {
	found := false
	// resetting last removed (since we haven't removed anything)
	t.lastRemove = ""

	len1 := len(t.curStr)
	len2 := len(what)

	for i := 0; i < len1 && i < len2; i++ {
		if t.curStr[i] != what[i] {
			if found {
				t.lastRemove = t.curStr[:i]
				t.curStr = t.curStr[len(t.lastRemove):] // equivalent to strings.TrimPrefix(t.curStr, t.lastRemove), but without HasPrefix call
			}
			return found
		}
		found = true
	}

	if !found {
		if len1 == len2 && t.curStr == what {
			// special case : "" not a subset of ""
			found = true
			t.lastRemove = t.curStr
			t.curStr = "" // would be  strings.TrimPrefix(t.curStr, t.curStr)
		}
	} else {
		if len1 > len2 {
			t.lastRemove = what
			t.curStr = t.curStr[len(what):] // equivalent to strings.TrimPrefix(t.curStr, what), but without HasPrefix call
		} else if len1 == len2 {
			t.lastRemove = t.curStr
			t.curStr = "" // would be strings.TrimPrefix(t.curStr, t.curStr)
		}
	}
	return found
}

func (t *Tree) search(target *Node) (interface{}, bool) {
	if target.isLeaf() {
		return target.leaf.value, true
	}

	for _, edge := range target.edges {
		if t.compare(edge.label) {
			return t.search(edge.child)
		}
	}

	return nil, false
}

func (t *Tree) starSearch(target *Node) (interface{}, bool) {
	if target.isLeaf() {
		return target.leaf.value, true
	}

	for _, edge := range target.edges {
		if edge.hasStar {
			// search key is empty and we're on the "/*" means we're looking for the last sibling of the edge
			if len(t.curStr) == 0 && edge.hasSlashStar {
				continue
			}
			// different kind of star... ("*" or "?*")
			t.compare(edge.label)

			// split by slashes so we can build a new key
			slashIdx := strings.IndexByte(t.curStr, slashByte)
			if slashIdx == -1 {
				// ok, we had one piece
				// looking for the question mark - might be handy to give up on this for speed
				index := strings.IndexByte(t.curStr, queByte) //index(t.curStr, que)
				if index > 0 {
					// collect param value
					t.params = append(t.params, t.curStr[index:])
					t.curStr = t.curStr[index:]
					// lookup question marks - down in the tree
					return t.starSearch(edge.child)
				}
				// don't have a question mark, but we have a star (continue) - looking for the last sibling edge
				if edge.hasQueStar && t.lastRemove != que {
					continue
				}

				// collect param value
				t.params = append(t.params, t.curStr)
				t.curStr = ""
				// we have a star, no question mark - looking for the node leaf
				return t.starSearch(edge.child)
			}

			// collect param value
			t.params = append(t.params, t.curStr[:slashIdx])
			// building a new key with the firstPart that we have
			t.curStr = t.curStr[slashIdx+1:]
			return t.starSearch(edge.child)
		}

		if t.compare(edge.label) {
			return t.starSearch(edge.child)
		}
	}

	return nil, false
}

func (t *Tree) Search(what string) (interface{}, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	// TODO : validation : if the string contains * or multiple ?
	t.curStr = what
	// reset the params collector
	t.params = make([]string, 0)
	if t.isStar {
		return t.starSearch(&t.root)
	}
	return t.search(&t.root)
}

func (t *Tree) insert(target *Node, edgeKey string, leafKey string, value interface{}) {
	// we've reached leaf
	if target.isLeaf() {
		if leafKey == target.leaf.key {
			// the same leafKey, update value
			target.leaf.value = value
		} else {
			// insert leaf key value as new child node of target
			// original leaf node, became another leaf of target
			target.createNode(edgeKey, leafKey, value)
			target.createNode("", target.leaf.key, target.leaf.value)
			target.leaf = nil
		}
		return
	}

	// second case - checking edges
	for _, edge := range target.edges {
		t.curStr = edgeKey
		if t.compare(edge.label) {
			if t.lastRemove == edge.label {
				// recurse
				t.insert(edge.child, t.curStr, leafKey, value)
				return
			}
			// else
			if t.lastRemove == "" {
				// switching between curStr and lastRemove so below code work both cases
				t.curStr, t.lastRemove = t.lastRemove, t.curStr
			}

			// using compare to create new node and separate two edges
			splitNode := target.createNodeWithEdges(t.lastRemove, edge.label)
			if splitNode == nil {
				panic("Unexpected error on creating new node and separating edges")
			}
			splitNode.createNode(t.curStr, leafKey, value)
			return
		}
	}

	// new edge with new leafKey on leaf node
	target.createNode(edgeKey, leafKey, value)

}

func (t *Tree) Insert(what string, value interface{}) {
	// leaf key and edge key are the same
	t.insert(&t.root, what, what, value)
}

func (t *Tree) searchLeaf(curNode, parNode *Node, curWhat, what string) (*Node, *Node, bool) {
	if curNode.isLeaf() {
		return curNode, parNode, curNode.leaf.key == what
	}

	for _, edge := range curNode.edges {
		t.curStr = curWhat
		if t.compare(edge.label) {
			return t.searchLeaf(edge.child, curNode, t.curStr, what)
		}
	}

	return nil, nil, false
}

func (t *Tree) SearchLeaf(what string) (*Node, *Node, bool) {
	return t.searchLeaf(&t.root, &t.root, what, what)
}

func (t *Tree) findParent(curNode, parNode, targetNode *Node) (*Node, bool) {
	if curNode.isLeaf() {
		return nil, false
	}

	if curNode == targetNode {
		return parNode, true
	}

	for _, edgeObj := range curNode.edges {
		if edgeObj.child == targetNode {
			return curNode, true
		}

		if pNode, find := t.findParent(edgeObj.child, curNode, targetNode); find {
			return pNode, true
		}
	}

	return nil, false
}

func (t *Tree) FindParent(target *Node) (*Node, bool) {
	return t.findParent(&t.root, &t.root, target)
}
