package skip

import (
	"math/rand"
	"time"
)

const (
	maxLevel = 32
	p        = 0.25
)

type List struct {
	header *Node
	level  int
	length int
	rnd    *rand.Rand
}

func New() *List {
	s := &List{
		header: newNode(maxLevel, 0, ""), // header with maxLevel arrays
		level:  1,
		rnd:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	return s
}

func (s *List) Header() *Node {
	return s.header
}

func (s *List) Length() int {
	return s.length
}

// randomLevel generates a random level for a new node
// with probability p for each level
// the lower the p, the less likely to have high levels
func (s *List) randomLevel() int {
	lvl := 1
	for lvl < maxLevel && s.rnd.Float64() < p {
		lvl++
	}
	return lvl
}

// Insert inserts a new (score, member) pair into the skiplist.
func (s *List) Insert(score float64, member string) *Node {
	insertionPoints, ranks := s.findInsertionPoints(score, member)

	lvl := s.pickLevelAndMaybeGrow(insertionPoints, ranks)

	x := s.createAndSpliceNode(score, member, lvl, insertionPoints, ranks)

	s.fixBackwardPointers(x, insertionPoints)

	s.length++
	return x
}

// findInsertionPoints walks the skiplist top-down to find where
// the new node should be inserted. Returns:
// - insertionPoints: the last node before the insertion point at each level
// - ranks: cumulative counts used for span calculation
func (s *List) findInsertionPoints(score float64, member string) ([]*Node, []int) {
	insertionPoints := make([]*Node, maxLevel)
	ranks := make([]int, maxLevel)

	x := s.header
	for i := s.level - 1; i >= 0; i-- {
		if i == s.level-1 {
			ranks[i] = 0
		} else {
			ranks[i] = ranks[i+1]
		}
		for x.nexts[i] != nil && less(x.nexts[i].score, x.nexts[i].member, score, member) {
			ranks[i] += x.span[i]
			x = x.nexts[i]
		}
		insertionPoints[i] = x
	}
	return insertionPoints, ranks
}

// pickLevelAndMaybeGrow chooses a random level and extends the skiplist if needed.
func (s *List) pickLevelAndMaybeGrow(insertionPoints []*Node, ranks []int) int {
	lvl := s.randomLevel()
	if lvl > s.level {
		for i := s.level; i < lvl; i++ {
			insertionPoints[i] = s.header
			insertionPoints[i].span[i] = s.length
			ranks[i] = 0
		}
		s.level = lvl
	}
	return lvl
}

// createAndSpliceNode creates the new node and links it at each level.
// Also updates span values for affected nodes.
func (s *List) createAndSpliceNode(score float64, member string, lvl int, insertionPoints []*Node, ranks []int) *Node {
	x := newNode(lvl, score, member)

	for i := 0; i < lvl; i++ {
		// forward pointers
		x.nexts[i] = insertionPoints[i].nexts[i]
		insertionPoints[i].nexts[i] = x

		// spans
		x.span[i] = insertionPoints[i].span[i] - (ranks[0] - ranks[i])
		insertionPoints[i].span[i] = (ranks[0] - ranks[i]) + 1
	}

	// untouched higher levels: increment span
	for i := lvl; i < s.level; i++ {
		insertionPoints[i].span[i]++
	}
	return x
}

// fixBackwardPointers updates level-0 backward links for doubly-linked traversal.
func (s *List) fixBackwardPointers(x *Node, insertionPoints []*Node) {
	if insertionPoints[0] == s.header {
		x.prev = nil
	} else {
		x.prev = insertionPoints[0]
	}
	if x.nexts[0] != nil {
		x.nexts[0].prev = x
	}
}

// Remove deletes the node with the given (score, member) from the skiplist.
// Returns true if the node was found and removed, false otherwise.
func (s *List) Remove(score float64, member string) bool {
	removalPoints := make([]*Node, maxLevel)
	x := s.findRemovalPoints(score, member, removalPoints)
	if x == nil {
		return false // not found
	}

	s.unlinkNode(x, removalPoints)
	s.fixBackwardPointersAfterRemove(x)
	s.length--

	// shrink level if needed
	for s.level > 1 && s.header.nexts[s.level-1] == nil {
		s.level--
	}
	return true
}

// findRemovalPoints locates the node to remove and records the nodes
// whose forward[i] pointers must be updated.
func (s *List) findRemovalPoints(score float64, member string, removalPoints []*Node) *Node {
	x := s.header
	for i := s.level - 1; i >= 0; i-- {
		for x.nexts[i] != nil && less(x.nexts[i].score, x.nexts[i].member, score, member) {
			x = x.nexts[i]
		}
		removalPoints[i] = x
	}

	// candidate node at level 0
	x = x.nexts[0]
	if x != nil && x.score == score && x.member == member {
		return x
	}
	return nil
}

// unlinkNode adjusts forward pointers and spans to bypass x.
func (s *List) unlinkNode(x *Node, removalPoints []*Node) {
	for i := 0; i < s.level; i++ {
		if removalPoints[i].nexts[i] == x {
			// bypass x
			removalPoints[i].span[i] += x.span[i] - 1
			removalPoints[i].nexts[i] = x.nexts[i]
		} else {
			// node not present at this level → just decrement span
			removalPoints[i].span[i]--
		}
	}
}

// fixBackwardPointersAfterRemove updates doubly-linked list pointers at level 0.
func (s *List) fixBackwardPointersAfterRemove(x *Node) {
	if x.nexts[0] != nil {
		x.nexts[0].prev = x.prev
	}
}

// GetRank returns the rank (0-based index) of the node with the given (score, member).
func (s *List) GetRank(score float64, member string) int {
	x := s.header
	rank := 0
	for i := s.level - 1; i >= 0; i-- {
		for x.nexts[i] != nil && less(x.nexts[i].score, x.nexts[i].member, score, member) {
			rank += x.span[i]
			x = x.nexts[i]
		}
	}

	x = x.nexts[0]
	if x != nil && x.score == score && x.member == member {
		return rank
	}

	return -1
}

// GetByRank returns the node at the given rank (0-based index).
func (s *List) GetByRank(rank int) *Node {
	if rank < 0 || rank >= s.length {
		return nil
	}
	x := s.header
	traversed := 0
	for i := s.level - 1; i >= 0; i-- {
		for x.nexts[i] != nil && traversed+x.span[i] <= rank {
			traversed += x.span[i]
			x = x.nexts[i]
		}
	}
	// x is the node before desired node
	return x.nexts[0]
}

// node comparison: first by score, then by member
func less(aScore float64, aMember string, bScore float64, bMember string) bool {
	if aScore < bScore {
		return true
	}

	if aScore > bScore {
		return false
	}

	return aMember < bMember
}
