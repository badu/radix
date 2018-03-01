package strRadix

import (
	"sort"
)

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
	n.edges.Add(Edge{
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
			n.edges.Replace(idx, Edge{
				label:  newKey,
				parent: n,
				child:  newNode,
			})
			// connect to original node
			remainKey := edgeKey[len(newKey):] //strings.TrimPrefix(edgeKey, newKey)
			newNode.edges.Add(Edge{
				label:  remainKey,
				parent: newNode,
				child:  oldNode,
			})
			return newNode
		}
	}
	return nil
}

func (e Edges) Len() int {
	return len(e)
}

func (e Edges) Less(i, j int) bool {
	return e[i].label > e[j].label
}

func (e Edges) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e Edges) Sort() {
	sort.Sort(e)
}

func (e *Edges) Replace(atIndex int, edge Edge) {
	// attention : first we're checking for star : causes malfunction because it's not a slice of pointers
	edge.setStars()
	// replace it in the slice
	(*e)[atIndex] = edge
	// we're always sorting in reverse, so stars are last siblings
	e.Sort()
}

func (e *Edges) Add(edge Edge) {
	// attention : first we're checking for star : causes malfunction because it's not a slice of pointers
	edge.setStars()
	// add it to the slice
	*e = append(*e, edge)
	// we're always sorting in reverse, so stars are last siblings
	e.Sort()
}

func (e *Edge) setStars() {
	e.hasQueStar = len(e.label) >= len(queStar) && e.label[0:len(queStar)] == queStar                       //strings.HasPrefix(e.label, queStar)
	e.hasSlashStar = len(e.label) >= len(slashStar) && e.label[0:len(slashStar)] == slashStar               //strings.HasPrefix(e.label, slashStar)
	e.hasStar = len(e.label) >= len(star) && e.label[0:len(star)] == star || e.hasQueStar || e.hasSlashStar //strings.HasPrefix(e.label, star) || e.hasQueStar || e.hasSlashStar
}
