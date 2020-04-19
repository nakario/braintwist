package brainfuck

import (
	"bytes"
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
