package sorted

import (
	"fmt"
	"testing"

	"github.com/codecrafters-io/redis-starter-go/pkg/skip"
)

func Test_range(t *testing.T) {
	list := skip.New()

	el1 := list.Insert(10, "v1")
	el2 := list.Insert(15, "v2")
	el3 := list.Insert(20, "v3")

	fmt.Printf("el1: key: %v, value: %v, score: %v\n", el1.Score(), el1.Member())
	fmt.Printf("el2: key: %v, value: %v, score: %v\n", el2.Score(), el2.Member())
	fmt.Printf("el3: key: %v, value: %v, score: %v\n", el3.Score(), el3.Member())
}
