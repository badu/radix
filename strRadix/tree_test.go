package strRadix

import (
	"testing"
)

type (
	TreeTester struct {
		Tree
	}
	routeAndValue struct {
		r string
		v int
	}
	routes []routeAndValue
)

const (
	Accept                  = "Accept"
	AcceptCharset           = "Accept-Charset"
	AcceptEncoding          = "Accept-Encoding"
	AcceptLanguage          = "Accept-Language"
	AcceptRanges            = "Accept-Ranges"
	Authorization           = "Authorization"
	CacheControl            = "Cache-Control"
	Cc                      = "Cc"
	Connection              = "Connection"
	ContentEncoding         = "Content-Encoding"
	ContentId               = "Content-Id"
	ContentLanguage         = "Content-Language"
	ContentLength           = "Content-Length"
	ContentRange            = "Content-Range"
	ContentTransferEncoding = "Content-Transfer-Encoding"
	ContentType             = "Content-Type"
	CookieHeader            = "Cookie"
	Date                    = "Date"
	DkimSignature           = "Dkim-Signature"
	Etag                    = "Etag"
	Expires                 = "Expires"
	Expect                  = "Expect"
	From                    = "From"
	Host                    = "Host"
	IfModifiedSince         = "If-Modified-Since"
	IfNoneMatch             = "If-None-Match"
	InReplyTo               = "In-Reply-To"
	LastModified            = "Last-Modified"
	Location                = "Location"
	MessageId               = "Message-Id"
	MimeVersion             = "Mime-Version"
	Pragma                  = "Pragma"
	Received                = "Received"
	Referer                 = "Referer"
	ReturnPath              = "Return-Path"
	ServerHeader            = "Server"
	SetCookieHeader         = "Set-Cookie"
	Subject                 = "Subject"
	TransferEncoding        = "Transfer-Encoding"
	To                      = "To"
	Trailer                 = "Trailer"
	UpgradeHeader           = "Upgrade"
	UserAgent               = "User-Agent"
	Via                     = "Via"
	XForwardedFor           = "X-Forwarded-For"
	XImforwards             = "X-Imforwards"
	XPoweredBy              = "X-Powered-By"
)

var (
	headers = []string{
		Accept,
		AcceptCharset,
		AcceptEncoding,
		AcceptLanguage,
		AcceptRanges,
		Authorization,
		CacheControl,
		Cc,
		Connection,
		ContentEncoding,
		ContentId,
		ContentLanguage,
		ContentLength,
		ContentRange,
		ContentTransferEncoding,
		ContentType,
		CookieHeader,
		Date,
		DkimSignature,
		Etag,
		Expires,
		Expect,
		From,
		Host,
		IfModifiedSince,
		IfNoneMatch,
		InReplyTo,
		LastModified,
		Location,
		MessageId,
		MimeVersion,
		Pragma,
		Received,
		Referer,
		ReturnPath,
		ServerHeader,
		SetCookieHeader,
		Subject,
		TransferEncoding,
		To,
		Trailer,
		UpgradeHeader,
		UserAgent,
		Via,
		XForwardedFor,
		XImforwards,
		XPoweredBy,
	}
)

//PrintTree: Print out current tree struct, it will using \t for tree level
func (t *TreeTester) PrintTree(currentNode *Node, treeLevel int) {
	if currentNode == nil {
		currentNode = &t.root
	}
	tabs := ""
	for i := 1; i < treeLevel; i++ {
		tabs = tabs + "\t"
	}

	if currentNode.isLeaf() {
		//Reach  the end point
		t.logger("%s[%d] Leaf key : %q value : %v\n", tabs, treeLevel, currentNode.leaf.key, currentNode.leaf.value)
		return
	}

	t.logger("%s[%d] Node has %d edges \n", tabs, treeLevel, len(currentNode.edges))
	for _, edge := range currentNode.edges {
		if edge.hasStar {
			t.logger("%s[%d] StarEdge [%q]\n", tabs, treeLevel, edge.label)
		} else {
			t.logger("%s[%d] NormalEdge [%q]\n", tabs, treeLevel, edge.label)
		}
		t.PrintTree(edge.child, treeLevel+1)
	}

	if treeLevel == 1 {
		t.logger("Tree printed.\n\n")
	}
}

func TestPrintTree(t *testing.T) {

	rootNode := Node{leaf: nil}

	cNode := Node{leaf: nil}
	lNode := Node{leaf: &leafNode{key: "company", value: 1}}
	rNode := Node{leaf: &leafNode{key: "comflict", value: 2}}

	rootEdge := Edge{label: "com"}
	rootEdge.parent = &rootNode
	rootEdge.child = &cNode
	rootNode.edges = append(rootNode.edges, rootEdge)

	lEdge := Edge{label: "pany"}
	lEdge.parent = &cNode
	lEdge.child = &lNode

	rEdge := Edge{label: "flict"}
	rEdge.parent = &cNode
	rEdge.child = &rNode

	cNode.edges = append(cNode.edges, lEdge)
	cNode.edges = append(cNode.edges, rEdge)

	t.Log("edges:", cNode.edges)
	rTree := TreeTester{Tree{logger: t.Logf}}
	rTree.root = rootNode

	rTree.PrintTree(nil, 1)
}

func TestNodeInsert(t *testing.T) {
	rTree := &TreeTester{Tree: Tree{logger: t.Logf}}

	rTree.root.createNode("keyAll", "keyAll", 1)
	rTree.root.createNode("open", "open", 2)
	rTree.PrintTree(nil, 1)
}

func TestTreeInsert(t *testing.T) {
	rTree := &TreeTester{Tree: Tree{logger: t.Logf}}
	rTree.Insert("test", 1)
	rTree.Insert("team", 2)

	if rTree.root.edges[0].label != "te" {
		t.Errorf("TreeInsert: Simple case failed, expect `te`, but get %s\n", rTree.root.edges[0].label)
	}

	rTree2 := &TreeTester{Tree: Tree{logger: t.Logf}}
	rTree2.Insert("main", 1)
	rTree2.Insert("mainly", 2)

	if rTree2.root.edges[0].label != "main" {
		t.Errorf("TreeInsert: Simple case failed, expect `main`, but get %s\n", rTree.root.edges[0].label)
	}
}

func TestLookup(t *testing.T) {
	rTree := &TreeTester{Tree: Tree{logger: t.Logf}}
	rTree.Insert("test", 1)
	rTree.Insert("team", 2)
	rTree.Insert("trobot", 3)
	rTree.Insert("apple", 4)
	rTree.Insert("app", 5)
	rTree.Insert("tesla", 6)

	ret, find := rTree.Search("team")
	if !find || ret != 2 {
		t.Errorf("Lookup failed, expect '2', but get %v", ret)
	}

	ret, find = rTree.Search("apple")
	if !find || ret != 4 {
		t.Errorf("Lookup failed, expect '4', but get %v", ret)
	}

	ret, find = rTree.Search("tesla")
	if !find || ret != 6 {
		t.Errorf("Lookup failed, expect '6', but get %v", ret)
	}

	ret, find = rTree.Search("app")
	if !find || ret != 5 {
		t.Errorf("Lookup failed, expect '5', but get %v", ret)
	}

	rTree.Insert("app", 7)
	rTree.PrintTree(nil, 1)
	ret, find = rTree.Search("app")
	t.Log(ret, find)
	if !find || ret != 7 {
		t.Errorf("Insert update lookup failed, expect '7', but get %v", ret)
	}
}

func TestLocateLeafNode(t *testing.T) {
	rTree := &TreeTester{Tree: Tree{logger: t.Logf}}
	rTree.Insert("test", 1)
	rTree.Insert("team", 2)
	rTree.Insert("trobot", 3)
	rTree.Insert("apple", 4)
	rTree.Insert("app", 5)
	rTree.Insert("tesla", 6)

	cNode, pNode, find := rTree.SearchLeaf("trobot")
	t.Log(cNode, pNode, find)

	cNode, pNode, find = rTree.SearchLeaf("trobota")
	t.Log(cNode, pNode, find)

	cNode, pNode, find = rTree.SearchLeaf("tesla")
	t.Log(cNode, pNode, find)
}

func TestFindParent(t *testing.T) {
	rTree := &TreeTester{Tree: Tree{logger: t.Logf}}
	rTree.Insert("test", 1)
	rTree.Insert("team", 2)
	rTree.Insert("trobot", 3)
	rTree.Insert("apple", 4)
	rTree.Insert("app", 5)
	rTree.Insert("tesla", 6)

	cNode, pNode, find := rTree.SearchLeaf("trobot")
	t.Log(cNode, pNode, find)
	cParent, cFind := rTree.FindParent(cNode)
	if cFind {
		t.Log(cParent.edges)
	} else {
		t.Errorf("Failed in find parentNode")
	}

	nextParent, ccFind := rTree.FindParent(cParent)
	if ccFind {
		t.Log(nextParent.edges)
	} else {
		t.Errorf("Failed in find parentNode")
	}

	pRoot, fRoot := rTree.FindParent(&rTree.root)
	if fRoot {
		if pRoot != &rTree.root {
			t.Errorf("Failed on find parent on root")
		}
		t.Log(pRoot.edges)
	} else {
		t.Errorf("Failed on find parent on root, cannot find it.")
	}
}

// unused
//PrintTree: Print out current tree struct, it will using \t for tree level
func (t *TreeTester) PrintNonRecursiveTree() {
	// we start at the root level
	currentNode := &t.root
	// building a queue to visit nodes
	queue := []*Node{}
	// while current node is not nil
	for currentNode != nil {

		if currentNode.isLeaf() {
			//Reach  the end point
			t.logger("Leaf key : %q value : %v\n", currentNode.leaf.key, currentNode.leaf.value)
		} else {
			// not a leaf node - has children
			t.logger("Node has %d edges \n", len(currentNode.edges))
			for _, edge := range currentNode.edges {
				if edge.hasStar {
					t.logger("StarEdge [%q]\n", edge.label)
				} else {
					t.logger("NormalEdge [%q]\n", edge.label)
				}
				queue = append(queue, edge.child)
			}
		}

		if len(queue) > 0 {
			// reading last in queue
			currentNode = queue[len(queue)-1]
			// deleting last in queue
			queue = queue[:len(queue)-1]
		} else {
			// exit condition
			currentNode = nil
		}

	}
}

func TestStarOneRoute(t *testing.T) {
	rTree := &TreeTester{Tree: Tree{logger: t.Logf, isStar: true}}
	rTree.Insert("*/*", 5555)
	// without the below insertion, works fine, but you won't be able to collect the parameters correctly
	rTree.Insert("*/never-matched", 10000)

	ret, find := rTree.Search("app/blah")
	if !find || ret != 5555 {
		t.Errorf("Lookup failed, expect '5555', but got %v", ret)
	} else {
		t.Log("Ok `app/blah` ", ret)
	}
	if len(rTree.params) > 0 {
		for idx, param := range rTree.params {
			t.Logf("\tparam %d = %q\n", idx, param)
		}
	}

	rTree.PrintTree(nil, 1)
	t.Log("You had one route.")
}

func TestStar(t *testing.T) {
	// Valid chars in URI -> https://tools.ietf.org/html/rfc3986#section-2
	// unwise      = "{" | "}" | "|" | "\" | "^" | "[" | "]" | "`"
	// reserved    = ";" | "/" | "?" | ":" | "@" | "&" | "=" | "+" | "$" | ","
	rTree := &TreeTester{Tree: Tree{logger: t.Logf, isStar: true}}

	def := routes{
		{
			"*", 11,
		}, {
			"test", 1,
		}, {
			"team", 2,
		}, {
			"trobot", 3,
		}, {
			"apple", 4,
		}, {
			"app", 5,
		}, {
			"app/blah", 555,
		}, {
			"app/blah/blah", 5555,
		}, {
			"app/blah/blah?*", 10000,
		}, {
			"tesla", 6,
		}, {
			"test/*", 12,
		}, {
			"test/*/*", 13,
		}, {
			"test/*/*?*", 14,
		}, {
			"tesla/copy/*?*", 202,
		}, {
			"tesla/calului/*?*", 220,
		}, {
			"tesla/calului/longer/*", 2020,
		}, {
			"tesla/calului/longer/|", 2021, // convention : for every star, we have to have "unmet" condition
		}, {
			"tesla/calului/larger/*", 2030,
		}, {
			"tesla/calului/larger/|", 2031, // convention : for every star, we have to have "unmet" condition
		}, {
			"tesla/*/*?*", 200,
		}, {
			"tesla/*/paste", 205,
		}, {
			"tesla/*/paste/*?*", 210,
		}, {
			"tșță", 1111,
		}, {
			"tșâță", 11111,
		},
	}

	for _, route := range def {
		rTree.Insert(route.r, route.v)
	}

	tests := []routeAndValue{
		{
			"app/blah/blah", 5555,
		}, {
			"app/blah/blah?filter=blah", 10000,
		}, {
			"app/blah", 555,
		}, {
			"tesla/copy/oops?search=blah", 202,
		}, {
			"test/457/doo?search=string", 14, //Note : test doesn't have defined siblings, that's why returns both `?search=string` and `search=string`
		}, {
			"test/123", 12,
		}, {
			"test/123/abc", 13,
		}, {
			"tesla/457/paste/oops?search=blah", 210,
		}, {
			"tesla/457/paste", 205,
		}, {
			"tesla/457/doo?search=string", 200,
		}, {
			"tesla/calului/longer/done", 2020,
		}, {
			"tesla/calului/larger/done", 2030,
		}, {
			"trobot", 3,
		}, {
			"wakarimasuka", 11,
		},
	}

	for _, test := range tests {
		t.Logf("Testint against %q\n", test.r)
		ret, find := rTree.Search(test.r)
		if !find || ret != test.v {
			t.Errorf("Lookup failed, expect '5555' on %q, but got %v", test.r, ret)
		} else {
			t.Logf("Ok %q - %d\n", test.r, ret)
		}
		if len(rTree.params) > 0 {
			for idx, param := range rTree.params {
				t.Logf("\tparam %d = %q\n", idx, param)
			}
		}
	}

	rTree.PrintTree(nil, 1)
	//t.Log("Test finished.")
}

func prepareTest(t *testing.T, tree *TreeTester) {
	for index, header := range headers {
		tree.Insert(header, index)
	}
	t.Log("Test prepared.")
}

func TestLookupHeaders(t *testing.T) {
	rTree := &TreeTester{Tree{logger: t.Logf}}
	prepareTest(t, rTree)

	ret, find := rTree.Search(Accept)
	if !find {
		t.Error("Lookup failed")
	} else {
		t.Logf("Found : %v", ret)
	}

	ret, find = rTree.Search(AcceptLanguage)
	if !find {
		t.Error("Lookup failed")
	} else {
		t.Logf("Found : %v", ret)
	}

	rTree.Insert("Foo-Header", 7)
	ret, find = rTree.Search("Foo-Header")
	if !find || ret != 7 {
		t.Errorf("Insert update lookup failed, expect '7', but get %v", ret)
	}
	t.Log(find, " found freshly inserted ", ret)

	rTree.PrintTree(nil, 1)

}
