package sorted

import (
	"fmt"
	"log/slog"
	"sync"

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
	mu   sync.RWMutex
	hmap map[string]*skip.Node
	list *skip.List
}

func (s *Set) Add(score float64, member string) (inserted bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	slog.Debug("ADD")

	el, exists := s.hmap[member]
	if exists {
		s.list.Remove(el.Score(), el.Member())
	} else {
		slog.Debug("INSERTED")
		inserted = true
	}

	el = s.list.Insert(score, member)
	s.hmap[member] = el

	return inserted
}

// Rank gives 0 for smallest score.
// If two have the same score, ranking is lexicological by the text of members: "aa" > "bb".
func (s *Set) Rank(member string) (int, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	el, ok := s.hmap[member]
	if !ok {
		return 0, false
	}

	return s.list.GetRank(el.Score(), el.Member()), true
}

// Range uses indexes (ranks), both inclusive.
func (s *Set) Range(from, to int) (result []string) {
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

func wrapNegativeIndex(v int, length int) int {
	if v < 0 {
		moveTimes := (-v)/length + 1
		fmt.Printf("moveTimes: %v\n", moveTimes)

		v += length * moveTimes
		v %= length
	}

	return v
}
