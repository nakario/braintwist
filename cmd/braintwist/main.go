package main

import (
	"flag"
	"fmt"
	"os"

	bt "github.com/nakario/braintwist"
)

var (
	memSize = flag.Int("memSize", 30000, "number of memory cells")
	eof = flag.Int("eof", 0, "value of EOF (0-255)")
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "Error: source code is not specified")
		fmt.Fprintln(os.Stderr, "Usage: braintwist src [option]")
		flag.Usage()
		os.Exit(1)
	}
	source := flag.Arg(0)
	f, err := os.Open(source)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: failed to open source code:", err)
		os.Exit(1)
	}
	prog, err := bt.Compile(f, bt.SetMemorySize(*memSize), bt.SetEOF(byte(*eof)))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: failed to compile program:", err)
		os.Exit(1)
	}
	if err := prog.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error: runtime error:", err)
		os.Exit(1)
	}
}
