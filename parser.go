package braintwist

import (
	"fmt"
	"strconv"
)

// Parser parses a Braintwist program.
type Parser struct {
	nodes []Node
	delay int
}

// NewParser returns a new Parser which takes tokens from lex.
func NewParser(lex *Lexer) *Parser {
	nodes := make([]Node, 0)
	for n := lex.Next(); n.Token != EOF; n = lex.Next() {
		nodes = append(nodes, n)
	}
	return &Parser{
		nodes: nodes,
		delay: 0,
	}
}

// Parse returns Generator constructed from a Braintwist program.
func (p *Parser) Parse() (Generator, error) {
	var generator Generator
	for _, n := range p.nodes {
		if n.Token == LF {
			p.delay++
			continue
		} else if n.Token != SEED {
			continue
		}
		seed, err := strconv.ParseUint(n.Value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse error (line %d: %s): %w", p.delay+1, n.Value, err)
		}
		if generator == nil {
			generator = NewGenerator(seed, p.delay)
		} else {
			generator = Xor(generator, NewGenerator(seed, p.delay))
		}
	}
	return generator, nil
}
