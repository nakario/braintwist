package brainfuck

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

// Token represents brainfuck's tokens.
type Token byte

func (t Token) String() string {
	return string([]byte{byte(t)})
}

// Tokens defined in brainfuck.
const (
	// Unknown is used to represent the zero value of Token type.
	Unknown Token = iota
	// EOF = End of File.
	EOF

	Inc   Token = '+'
	Dec   Token = '-'
	Next  Token = '>'
	Prev  Token = '<'
	Read  Token = ','
	Write Token = '.'
	Open  Token = '['
	Close Token = ']'
)

var validTokens = []Token{Inc, Dec, Next, Prev, Read, Write, Open, Close}

// ErrBOT is an error returned by Close operation when it can't find
// any Open token corresponding to it.
var ErrBOT = errors.New("beggining of tokens")

// Lexer reads Tokens from a source code.
type Lexer interface {
	Next() Token
}

type lexer struct {
	r io.Reader
}

// NewLexer returns a new Lexer which reads tokens from r.
func NewLexer(r io.Reader) Lexer {
	return &lexer{r: r}
}

// Next returns the next token.
// It ignores all characters other than +-<>,.[]
func (l lexer) Next() Token {
	for {
		c := make([]byte, 1)
		if _, err := l.r.Read(c); err != nil {
			if errors.Is(err, io.EOF) {
				return EOF
			}
			panic(fmt.Errorf("brainfuck lexer: unexpected error while reading source code: %w", err))
		}
		for _, vt := range validTokens {
			if c[0] == byte(vt) {
				return vt
			}
		}
		// All characters other than validTokens are comments.
	}
}

type tokenBuffer struct {
	buffer []Token
	lex    Lexer
	pos    int
}

func newTokenBuffer(l Lexer) *tokenBuffer {
	return &tokenBuffer{
		buffer: make([]Token, 0),
		lex:    l,
		pos:    0,
	}
}

func (tb *tokenBuffer) next() (Token, error) {
	if tb.pos < len(tb.buffer) {
		t := tb.buffer[tb.pos]
		tb.pos++
		return t, nil
	}
	t := tb.lex.Next()
	if t == EOF {
		return EOF, io.EOF
	}
	tb.buffer = append(tb.buffer, t)
	tb.pos++
	return t, nil
}

func (tb *tokenBuffer) prev() (Token, error) {
	if tb.pos <= 1 {
		return Unknown, ErrBOT
	}
	tb.pos--
	return tb.buffer[tb.pos-1], nil
}

// Option is a optional settings for Interpreter.
type Option func(p *Interpreter)

// SetMemorySize sets the number of memory cells.
func SetMemorySize(size int) Option {
	return Option(func(p *Interpreter) {
		p.mem = make([]byte, size)
	})
}

// SetInput sets Interpreter's input.
func SetInput(input io.Reader) Option {
	return Option(func(p *Interpreter) {
		p.reader = bufio.NewReader(input)
	})
}

// SetOutput sets Interpreter's output.
func SetOutput(output io.Writer) Option {
	return Option(func(p *Interpreter) {
		p.writer = output
	})
}

// Interpreter is a runtime for a brainfuck program.
type Interpreter struct {
	tokens *tokenBuffer
	mem    []byte
	ptr    int
	reader *bufio.Reader
	writer io.Writer
}

// NewInterpreter returns Interpreter running on tokens from Lexer.
// The default memory size is 30000.
// The default i/o is os.Stdin and os.Stdout.
func NewInterpreter(l Lexer, options ...Option) *Interpreter {
	p := &Interpreter{
		tokens: newTokenBuffer(l),
		ptr:    0,
	}
	for _, o := range options {
		o(p)
	}
	if p.mem == nil {
		SetMemorySize(30000)(p)
	}
	if p.reader == nil {
		SetInput(os.Stdin)(p)
	}
	if p.writer == nil {
		SetOutput(os.Stdout)(p)
	}
	return p
}

// Run executes all operations until it reaches the end of
// the program or an error occurs.
func (p *Interpreter) Run() error {
	for {
		finished, err := p.Step()
		if err != nil {
			return err
		}
		if finished {
			break
		}
	}
	return nil
}

// Step executes one operation.
func (p *Interpreter) Step() (finished bool, err error) {
	defer func() {
		if rerr := recover(); rerr != nil {
			// All panics recovered here should be "index out of range"
			err = errors.New("invalid cell address")
		}
	}()
	t, err := p.tokens.next()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return true, nil
		}
		return false, err
	}
	switch t {
	case Inc:
		p.inc()
	case Dec:
		p.dec()
	case Next:
		p.next()
	case Prev:
		p.prev()
	case Read:
		if err := p.read(); err != nil {
			return false, err
		}
	case Write:
		if err := p.write(); err != nil {
			return false, err
		}
	case Open:
		if err := p.open(); err != nil {
			return false, err
		}
	case Close:
		if err := p.close(); err != nil {
			if errors.Is(err, ErrBOT) {
				return true, nil
			}
			return false, err
		}
	}
	return false, nil
}

func (p *Interpreter) inc() {
	p.mem[p.ptr]++
}

func (p *Interpreter) dec() {
	p.mem[p.ptr]--
}

func (p *Interpreter) next() {
	p.ptr++
}

func (p *Interpreter) prev() {
	p.ptr--
}

func (p *Interpreter) read() error {
	b, err := p.reader.ReadByte()
	if err != nil {
		if errors.Is(err, io.EOF) {
			b = 0
		} else {
			return err
		}
	}
	p.mem[p.ptr] = b
	return nil
}

func (p *Interpreter) write() error {
	if _, err := p.writer.Write([]byte{p.mem[p.ptr]}); err != nil {
		return err
	}
	return nil
}

func (p *Interpreter) open() error {
	if p.mem[p.ptr] != 0 {
		return nil
	}
	nest := 0
	for {
		t, err := p.tokens.next()
		if err != nil {
			return err
		}
		if t == Open {
			nest++
		} else if t == Close {
			if nest == 0 {
				break
			}
			nest--
		}
	}
	return nil
}

func (p *Interpreter) close() error {
	if p.mem[p.ptr] == 0 {
		return nil
	}
	nest := 0
	for {
		t, err := p.tokens.prev()
		if err != nil {
			return err
		}
		if t == Close {
			nest++
		} else if t == Open {
			if nest == 0 {
				break
			}
			nest--
		}
	}
	return nil
}
