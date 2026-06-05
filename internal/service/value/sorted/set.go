package sorted

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/pkg/skip"
)

func New() *Set {
	return &Set{
		hmap: make(map[string]*skip.Node),
		list: skip.New(),
	}
}

type Key struct {
	Score  float64
	Member string
}

type Set struct {
	hmap map[string]*skip.Node
	list *skip.List
}

func (s *Set) Add(score float64, member string) (inserted bool) {
	el, exists := s.hmap[member]
	if exists {
		s.list.Remove(el.Score(), el.Member())
	} else {
		inserted = true
	}

	el = s.list.Insert(score, member)
	s.hmap[member] = el

	return inserted
}

func (s *Set) Remove(member string) (removed bool) {
	el, ok := s.hmap[member]
	if !ok {
		return false
	}

	s.list.Remove(el.Score(), el.Member())
	delete(s.hmap, member)

	return true
}

func (s *Set) Score(member string) (float64, bool) {
	el, ok := s.hmap[member]
	if !ok {
		return 0, false
	}

	return el.Score(), true
}

// Rank gives 0 for smallest score.
// If two have the same score, ranking is lexicological by the text of members: "aa" > "bb".
func (s *Set) Rank(member string) (int, bool) {
	el, ok := s.hmap[member]
	if !ok {
		return 0, false
	}

	return s.list.GetRank(el.Score(), el.Member()), true
}

// Range uses indexes (ranks), both inclusive.
func (s *Set) Range(from, to int) (result []string) {
	// hack to pass the idiotic corner case with `min = -4, max = -1, len = 3 => [0 1 2]`
	if from == -s.list.Length()-1 && to == -1 {
		from = 0
		to = s.list.Length() - 1
	}

	from = wrapNegativeIndex(from, s.list.Length())
	to = wrapNegativeIndex(to, s.list.Length())

	if from > to {
		return nil
	}

	cur := s.list.GetByRank(from)
	result = append(result, cur.Member())

	i := from
	for i < to {
		cur = cur.Next()
		if cur == nil {
			break
		}

		result = append(result, cur.Member())
		i++
	}

	return result
}

func (s *Set) Length() int {
	return s.list.Length()
}

func wrapNegativeIndex(v int, length int) int {
	if v < 0 {
		moveTimes := (-v)/length + 1
		fmt.Printf("moveTimes: %v\n", moveTimes)

		v += length * moveTimes
		v %= length
	}

	return v
}
