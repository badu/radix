package strRadix

import "sync"

const (
	queStar   = "?*"
	slashStar = "/*"
	star      = "*"
	slash     = "/"
	que       = "?"
)

type (
	Tree struct {
		mu         sync.Mutex // mutex to guard Search
		root       Node
		isStar     bool   // marks searching for star or normal radix (no special cases)
		lastRemove string // temporary string to avoid allocations
		curStr     string // temporary string to avoid allocations
		logger     func(format string, args ...interface{})
		params     []string
	}

	Edge struct {
		hasStar      bool
		label        string
		parent       *Node
		child        *Node
		hasSlashStar bool
		hasQueStar   bool
	}

	Edges []Edge

	leafNode struct {
		key   string
		value interface{}
	}

	Node struct {
		leaf  *leafNode
		edges Edges
	}
)
