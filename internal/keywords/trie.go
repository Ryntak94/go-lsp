package keywords

type Node struct {
	Value    string
	Children map[string]*Node
	Complete bool
	Keywords KeywordMap
}

type Trie struct {
	Root *Node
}

type KeywordMap map[Keyword]Keyword

func (km KeywordMap) add(keyword Keyword) {
	km[keyword] = keyword
}

func CreateTrie() Trie {
	n := Node{Children: map[string]*Node{}}
	return Trie{Root: &n}
}

func (t *Trie) AddWord(keyword Keyword) {
	var currNode *Node = t.Root
	for _, c := range keyword.Name() {
		char := string(c)
		currNode = currNode.addChild(char)
	}
	currNode.Complete = true
	currNode.Keywords.add(keyword)
}

func (t *Trie) FindWords() KeywordMap {
	keywords := KeywordMap{}
	for _, child := range t.Root.Children {
		child.findWords("", keywords)
	}
	return keywords
}

func createNode(v string) *Node {
	n := Node{Value: v, Children: map[string]*Node{}, Keywords: KeywordMap{}}
	return &n
}

func (n *Node) getValue() string {
	return n.Value
}

func (n *Node) addChild(char string) *Node {
	if child := n.getChild(char); child != nil {
		return child
	}
	child := createNode(char)
	n.Children[char] = child
	return child
}

func (n *Node) getChild(k string) *Node {
	return n.Children[k]
}

func (n *Node) findWords(subString string, keywords KeywordMap) {
	currSubString := subString + n.getValue()
	if n.Complete == true {
		for _, keyword := range n.Keywords {
			keywords.add(keyword)
		}
	}
	for _, child := range n.Children {
		child.findWords(currSubString, keywords)
	}
}
