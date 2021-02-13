package brainfuck

import (
	"bytes"
	"errors"
	"testing"
)

func TestLexer(t *testing.T) {
	src := ",. +[-<\n>] a"
	r := bytes.NewBufferString(src)
	l := NewLexer(r)
	tokens := make([]Token, 0)
	for t := l.Next(); t != EOF; t = l.Next() {
		tokens = append(tokens, t)
	}
	expects := []Token{
		Read,
		Write,
		Inc,
		Open,
		Dec,
		Prev,
		Next,
		Close,
	}
	for i, e := range expects {
		if tokens[i] != e {
			t.Errorf(
				"Token[%d]: %s != expected: %s",
				i,
				string([]byte{byte(tokens[i])}),
				string([]byte{byte(e)}),
			)
		}
	}
}

func TestHelloWorld(t *testing.T) {
	src := "++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>>.<-.<.+++.------.--------.>>+.>++."
	l := NewLexer(bytes.NewBufferString(src))
	out := new(bytes.Buffer)
	prog := NewInterpreter(l, SetOutput(out))
	if err := prog.Run(); err != nil {
		t.Error(err)
	}
	expected := "Hello World!\n"
	output := out.String()
	if output != expected {
		t.Errorf("program output(%s) != expected(%s)", output, expected)
	}
}

func TestCat(t *testing.T) {
	src := ",[.,]"
	text := "Test 123"
	input := bytes.NewBufferString(text)
	output := new(bytes.Buffer)
	l := NewLexer(bytes.NewBufferString(src))
	prog := NewInterpreter(l, SetInput(input), SetOutput(output))
	if err := prog.Run(); err != nil {
		t.Error(err)
	}
	outputStr := output.String()
	if outputStr != text {
		t.Errorf("program output(%s) != expected(%s)", outputStr, text)
	}
}

func TestCatEOF255(t *testing.T) {
	src := ",+[-.,+]"
	text := "Test 123"
	input := bytes.NewBufferString(text)
	output := new(bytes.Buffer)
	l := NewLexer(bytes.NewBufferString(src))
	prog := NewInterpreter(l, SetInput(input), SetOutput(output), SetEOF(255))
	if err := prog.Run(); err != nil {
		t.Error(err)
	}
	outputStr := output.String()
	if outputStr != text {
		t.Errorf("program output(%s) != expected(%s)", outputStr, text)
	}
}

func TestNegativePointer(t *testing.T) {
	src := "<<>>,."
	text := "Test 123"
	input := bytes.NewBufferString(text)
	output := new(bytes.Buffer)
	l := NewLexer(bytes.NewBufferString(src))
	prog := NewInterpreter(l, SetInput(input), SetOutput(output))
	if err := prog.Run(); err != nil {
		t.Error(err)
	}
	outputStr := output.String()
	if outputStr != "T" {
		t.Errorf("program output(%s) != expected(%s)", outputStr, "T")
	}
}

func TestNegativeDereference(t *testing.T) {
	src := "<<+"
	l := NewLexer(bytes.NewBufferString(src))
	prog := NewInterpreter(l)
	if err := prog.Run(); !errors.Is(err, ErrInvalidAddress) {
		t.Errorf("caused error(%v) != expected(%v)", err, ErrInvalidAddress)
	}
}

func TestOutOfRangeDereference(t *testing.T) {
	src := ">>>+"
	l := NewLexer(bytes.NewBufferString(src))
	prog := NewInterpreter(l, SetMemorySize(2))
	if err := prog.Run(); !errors.Is(err, ErrInvalidAddress) {
		t.Errorf("caused error(%v) != expected(%v)", err, ErrInvalidAddress)
	}
}

func TestBOT(t *testing.T) {
	src := "[,.>]+[,+.>]],++.>+]"
	l := NewLexer(bytes.NewBufferString(src))
	input := bytes.NewBufferString("ADKP")
	output := new(bytes.Buffer)
	prog := NewInterpreter(l, SetInput(input), SetOutput(output))
	if err := prog.Run(); err != nil {
		t.Error(err)
	}
	expected := "BF"
	if output.String() != expected {
		t.Errorf("program output(%s) != expected(%s)", output.String(), expected)
	}
}
