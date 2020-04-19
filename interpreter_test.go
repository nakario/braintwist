package braintwist

import (
	"bytes"
	"testing"
)

func TestCat(t *testing.T) {
	src := "848406"
	text := "Test 123"
	input := bytes.NewBufferString(text)
	output := new(bytes.Buffer)
	i, err := Compile(bytes.NewBufferString(src), SetInput(input), SetOutput(output))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := i.Run(); err != nil {
		t.Error(err)
		t.FailNow()
	}
	v := output.String()
	if v != text {
		t.Errorf("output (%s) != expected (%s)", v, text)
	}
}

func TestHello(t *testing.T) {
	src := HelloWorld
	output := new(bytes.Buffer)
	i, err := Compile(bytes.NewBufferString(src), SetOutput(output))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := i.Run(); err != nil {
		t.Error(err)
		t.FailNow()
	}
	v := output.String()
	text := "Hello World!\n"
	if v != text {
		t.Errorf("output (%s) != expected (%s)", v, text)
	}
}

const HelloWorld = `







956



9

365


371


484


223


419


974





843


459


302



97


448


462


140



54


431


860


20


67

734


643



381


20


134



9


53


847


255




654


3`
