package braintwist

import (
	"math/rand"

	"github.com/seehuhn/mt19937"
)

func xor(v1, v2 uint64) uint64 {
	return v1 ^ v2
}

// Xor returns a new Generator whose elements are
// xor of elements of g1 and g2.
func Xor(g1, g2 Generator) Generator {
	return newUnionGenerator(g1, g2, xor)
}

// Generator is a pseudo random number generator with some
// delay of starting generation.
type Generator interface {
	Generate() uint64
}

type generator struct {
	seed  uint64
	delay int
	prng  *rand.Rand
}

// NewGenerator returns a Generator which is based on 64 bit Mersenne twister
// with delay of starting generation.
func NewGenerator(seed uint64, delay int) Generator {
	prng := rand.New(mt19937.New())
	prng.Seed(int64(seed))
	return &generator{
		seed:  seed,
		delay: delay,
		prng:  prng,
	}
}

func (g *generator) Generate() uint64 {
	if g.delay > 0 {
		g.delay--
		return 0
	}
	return g.prng.Uint64()
}

type unionGenerator struct {
	g1, g2 Generator
	op     func(v1, v2 uint64) uint64
}

func newUnionGenerator(g1, g2 Generator, op func(v1, v2 uint64) uint64) Generator {
	return &unionGenerator{g1: g1, g2: g2, op: op}
}

func (g *unionGenerator) Generate() uint64 {
	return g.op(g.g1.Generate(), g.g2.Generate())
}
