package strRadix

func (n *Node) isLeaf() bool {
	return n.leaf != nil && len(n.edges) == 0
}

func (n *Node) createNode(edgeKey, leafKey string, value interface{}) {
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

func (n *Node) createNodeWithEdges(newKey string, edgeKey string) *Node {
	if n.isLeaf() {
		//node is leaf node could not split, return nil
		return nil
	}

	for idx, edge := range n.edges {
		if edge.label == edgeKey {
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
			// connect to original node
			remainKey := edgeKey[len(newKey):] //strings.TrimPrefix(edgeKey, newKey)
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
	e.hasQueStar = len(e.label) >= 2 && e.label[:2] == queStar                             //strings.HasPrefix(e.label, queStar)
	e.hasSlashStar = len(e.label) >= 2 && e.label[:2] == slashStar                         //strings.HasPrefix(e.label, slashStar)
	e.hasStar = len(e.label) >= 1 && e.label[:1] == star || e.hasQueStar || e.hasSlashStar //strings.HasPrefix(e.label, star) || e.hasQueStar || e.hasSlashStar
}
