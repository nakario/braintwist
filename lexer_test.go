package braintwist

import (
	"strings"
	"testing"
)

func TestLexer(t *testing.T) {
	src := strings.NewReader("123 a456   789.012")
	l := NewLexer(src)
	nodes := make([]Node, 0)
	for n := l.Next(); n.Token != EOF; n = l.Next() {
		nodes = append(nodes, n)
	}
	expects := []Node{
		{SEED, "123", 1},
		{SEED, "456", 6},
		{SEED, "789", 12},
		{SEED, "012", 16},
	}
	for i, e := range expects {
		if nodes[i] != e {
			t.Errorf("Node[%d]: %v != expected: %v\n", i, nodes[i], e)
		}
	}
}
