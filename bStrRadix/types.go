package strRadix

import "sync"

var (
	zeroSlice      = make([]byte, 0)
	queStar        = []byte("?*")
	slashStar      = []byte("/*")
	star           = []byte("*")
	slashByte byte = '/'
	que            = []byte("?")
	queByte   byte = '?'
)

type (
	Tree struct {
		mu         sync.Mutex // mutex to guard Search
		root       Node
		IsStar     bool   // marks searching for star or normal radix (no special cases)
		lastRemove []byte // temporary string to avoid allocations
		curStr     []byte // temporary string to avoid allocations
		logger     func(format string, args ...interface{})
		params     [][]byte
	}

	Edge struct {
		label        []byte
		parent       *Node
		child        *Node
		hasStar      bool
		hasSlashStar bool
		hasQueStar   bool
	}

	Edges []Edge

	leafNode struct {
		key   []byte
		value interface{}
	}

	Node struct {
		leaf  *leafNode
		edges Edges
	}
)
