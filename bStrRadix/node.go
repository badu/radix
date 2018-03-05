package strRadix

func (n *Node) isLeaf() bool {
	return n.leaf != nil && len(n.edges) == 0
}

func (n *Node) createNode(edgeKey, leafKey []byte, value interface{}) {
	leafNode := &Node{
		leaf: &leafNode{
			key:   leafKey,
			value: value,
		},
	}
	n.edges.add(Edge{
		label:  edgeKey,
		parent: n,
		child:  leafNode,
	})
}

func (n *Node) createNodeWithEdges(newKey []byte, edgeKey []byte) *Node {
	if n.isLeaf() {
		//node is leaf node could not split, return nil
		return nil
	}
	var remainKey []byte
	newKeySize := len(newKey)
	edgeKeySize := len(edgeKey)
	for idx, edge := range n.edges {
		if equal(edge.label, edgeKey) {
			// backup for split
			oldNode := edge.child
			// createNodeWithEdges split node
			newNode := &Node{}
			// replace current edge with a new one
			n.edges.replace(idx, Edge{
				label:  newKey,
				parent: n,
				child:  newNode,
			})
			// connect to original node - read as `remainKey := bytes.TrimPrefix(edgeKey, newKey)`
			if edgeKeySize >= newKeySize && equal(edgeKey[:newKeySize], newKey) {
				remainKey = edgeKey[newKeySize:]
			} else {
				remainKey = edgeKey
			}
			newNode.edges.add(Edge{
				label:  remainKey,
				parent: newNode,
				child:  oldNode,
			})
			return newNode
		}
	}
	return nil
}

func (e *Edge) setStars() {
	e.hasQueStar = len(e.label) >= 2 && equal(e.label[:2], queStar)
	e.hasSlashStar = len(e.label) >= 2 && equal(e.label[:2], slashStar)
	e.hasStar = len(e.label) >= 1 && equal(e.label[:1], star) || e.hasQueStar || e.hasSlashStar
}
