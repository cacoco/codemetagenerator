package utils

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestFilter1(t *testing.T) {
	g := gomega.NewWithT(t)

	input := []string{"a", "b", "c", "d", "e", "f"}
	predicate := func(s string) bool {
		return s < "d"
	}

	result := Filter(input, predicate)
	g.Expect(result).To(gomega.Equal([]string{"a", "b", "c"}))
}

func TestFilter2(t *testing.T) {
	g := gomega.NewWithT(t)

	input := []int{1, 2, 3, 4, 5, 6}
	predicate := func(i int) bool {
		return i%2 == 0
	}

	result := Filter(input, predicate)
	g.Expect(result).To(gomega.Equal([]int{2, 4, 6}))
}

func TestFilter3(t *testing.T) {
	g := gomega.NewWithT(t)

	input := []string{"a", "b", "", "d", "e", "f"}
	predicate := func(s string) bool {
		return len(s) > 0
	}

	result := Filter(input, predicate)
	g.Expect(result).To(gomega.Equal([]string{"a", "b", "d", "e", "f"}))
}

func TestFilter4(t *testing.T) {
	g := gomega.NewWithT(t)

	input := []any{"hello", 1, "world", 2, "!"}
	predicate := func(a any) bool {
		switch a.(type) {
		case string:
			return true
		default:
			return false
		}
	}

	result := Filter(input, predicate)
	g.Expect(result).To(gomega.Equal([]any{"hello", "world", "!"}))
}
