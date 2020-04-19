package braintwist

import (
	"testing"
)

func TestGenerator(t *testing.T) {
	g := NewGenerator(0, 0)
	values := make([]uint64, 3)
	for i := 0; i < 3; i++ {
		values[i] = g.Generate()
	}
	expects := []uint64{
		2947667278772165694,
		18301848765998365067,
		729919693006235833,
	}
	for i, e := range expects {
		if values[i] != e {
			t.Errorf("generated value[%d]: %d != expected value: %d", i, values[i], e)
		}
	}
}

func TestMersenneTwister(t *testing.T) {
	g := NewGenerator(5489, 0)
	for i := 0; i < 9999; i++ {
		g.Generate()
	}
	v := g.Generate()
	var expect uint64 = 9981545732273789042
	if v != expect {
		t.Errorf("the 10000th generated value (%d) != %d", v, expect)
	}
}

func TestGeneratorDelay(t *testing.T) {
	g := NewGenerator(5489, 3)
	for i := 0; i < 3; i++ {
		if v := g.Generate(); v != 0 {
			t.Errorf("non-zero value (%d) was generated", v)
		}
	}
	for i := 0; i < 9999; i++ {
		g.Generate()
	}
	v := g.Generate()
	var expect uint64 = 9981545732273789042
	if v != expect {
		t.Errorf("the 10000th generated value with delay (%d) != %d", v, expect)
	}
}

func TestUnionGenerator(t *testing.T) {
	g1 := NewGenerator(0, 1)
	g2 := NewGenerator(1, 0)
	union := Xor(g1, g2)
	values := make([]uint64, 3)
	for i := range values {
		values[i] = union.Generate()
	}
	expects := []uint64{
		0 ^ 2469588189546311528,
		2947667278772165694 ^ 2516265689700432462,
		18301848765998365067 ^ 8323445853463659930,
	}
	for i, e := range expects {
		if values[i] != e {
			t.Errorf("generated value[%d]: %d != expected value: %d", i, values[i], e)
		}
	}
}

func TestNestedUnionGenerator(t *testing.T) {
	g1 := NewGenerator(0, 1)
	g2 := NewGenerator(1, 0)
	g3 := NewGenerator(2, 0)
	union := Xor(Xor(g1, g2), g3)
	values := make([]uint64, 3)
	for i := range values {
		values[i] = union.Generate()
	}
	expects := []uint64{
		0 ^ 2469588189546311528 ^ 16668552215174154828,
		2947667278772165694 ^ 2516265689700432462 ^ 15684088468973760345,
		18301848765998365067 ^ 8323445853463659930 ^ 14458935525009338917,
	}
	for i, e := range expects {
		if values[i] != e {
			t.Errorf("generated value[%d]: %d != expected value: %d", i, values[i], e)
		}
	}
}
