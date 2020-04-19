package braintwist

import (
	"io"

	bf "github.com/nakario/braintwist/brainfuck"
)

// Interpreter controlls evaluation of a braintwist program.
type Interpreter struct {
	bfi *bf.Interpreter
}

var bfTokens = []bf.Token{
	bf.Inc, bf.Dec, bf.Next, bf.Prev,
	bf.Read, bf.Write, bf.Open, bf.Close,
}

type bfLexer struct {
	gen Generator
}

// Option is an optional setting for the internal brainfuck interpreter.
type Option = bf.Option

// These options come from github.com/nakario/braintwist/brainfuck
var (
	SetMemorySize = bf.SetMemorySize
	SetInput      = bf.SetInput
	SetOutput     = bf.SetOutput
)

func (l bfLexer) Next() bf.Token {
	return bfTokens[l.gen.Generate()&7] // 7 masks the bottom 3 bits
}

// NewInterpreter returns a new Interpreter which runs on a brainfuck program
// generated from Generator.
// The default options are equal to github.com/nakario/braintwist/brainfuck#Interpreter
func NewInterpreter(g Generator, options ...Option) *Interpreter {
	return &Interpreter{
		bfi: bf.NewInterpreter(&bfLexer{gen: g}, options...),
	}
}

// Compile is an utility function which returns Interpreter directly from a source code.
func Compile(src io.Reader, options ...Option) (*Interpreter, error) {
	l := NewLexer(src)
	p := NewParser(l)
	g, err := p.Parse()
	if err != nil {
		return nil, err
	}
	return NewInterpreter(g, options...), nil
}

// Run executes the program until it finishes.
func (i *Interpreter) Run() error {
	return i.bfi.Run()
}
