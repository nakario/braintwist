package braintwist

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Token is the type of tokens in Braintwist.
type Token byte

// Tokens defined in Braintwist.
const (
	// ILLEGAL is the zero-value of Token type.
	ILLEGAL Token = iota
	// EOF = End of File
	EOF
	// LF = Line Feed
	LF

	// SEED represents numerals
	SEED
)

// Node represents token information.
type Node struct {
	Token Token
	Value string
	Pos   int
}

// Lexer reads tokens from a source code.
type Lexer struct {
	reader *bufio.Reader
	pos    int
}

// NewLexer returns Lexer which reads from r.
func NewLexer(r io.Reader) *Lexer {
	return &Lexer{
		reader: bufio.NewReader(r),
		pos:    0,
	}
}

// Next returns the next token as Node.
func (l *Lexer) Next() Node {
	var numeral strings.Builder
	for {
		l.pos++
		c, err := l.reader.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if numeral.Len() > 0 {
					l.pos--
					num := numeral.String()
					numeral.Reset()
					return Node{
						Token: SEED,
						Value: num,
						Pos:   l.pos - len(num) + 1,
					}
				}
				return Node{
					Token: EOF,
					Value: "",
					Pos:   l.pos,
				}
			}
			panic(fmt.Errorf("unexpected error while reading source code: %w", err))
		}
		if '0' <= c && c <= '9' {
			numeral.WriteByte(c)
		} else if numeral.Len() > 0 {
			l.pos--
			_ = l.reader.UnreadByte()
			num := numeral.String()
			numeral.Reset()
			return Node{
				Token: SEED,
				Value: num,
				Pos:   l.pos - len(num) + 1,
			}
		}
		if c == '\n' {
			return Node{
				Token: LF,
				Value: "\n",
				Pos:   l.pos,
			}
		}
	}
}
