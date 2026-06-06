package skip

type Node struct {
	score  float64
	member string
	nexts  []*Node // nexts[i] points to next node at level i
	span   []int   // span[i] = #nodes between this node and nexts[i]
	prev   *Node   // level 0 backward pointer for reverse iteration
}

func newNode(level int, score float64, member string) *Node {
	return &Node{
		score:  score,
		member: member,
		nexts:  make([]*Node, level),
		span:   make([]int, level),
		prev:   nil,
	}
}

func (n *Node) Prev() *Node {
	return n.prev
}

func (n *Node) Next() *Node {
	return n.NextByLevel(0)
}

func (n *Node) NextByLevel(level int) *Node {
	if level < 0 || level >= len(n.nexts) {
		return nil
	}
	return n.nexts[level]
}

func (n *Node) Member() string {
	return n.member
}

func (n *Node) Score() float64 {
	return n.score
}
