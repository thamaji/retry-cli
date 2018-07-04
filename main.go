package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/thamaji/ioutils"
	"golang.org/x/crypto/ssh/terminal"
)

const Version = "v1.0.0"

func main() {
	var help, version bool
	var n, delay int
	var silent bool

	flag.BoolVar(&help, "h", false, "show help")
	flag.BoolVar(&version, "v", false, "show version")

	flag.IntVar(&n, "n", 3, "set number of max retry count")
	flag.IntVar(&delay, "d", 10, "set delay seconds")
	flag.BoolVar(&silent, "s", false, "silent mode")

	flag.Usage = func() {
		fmt.Println()
		fmt.Println("Usage: " + os.Args[0] + " [OPTIONS] COMMAND [ARG...]")
		fmt.Println()
		fmt.Println("Retry command if finished with non-zero")
		fmt.Println()
		fmt.Println("Options:")
		flag.CommandLine.PrintDefaults()
		fmt.Println()
	}

	flag.Parse()

	args := flag.Args()

	if help {
		flag.Usage()
		return
	}

	if version {
		fmt.Println(Version)
		return
	}

	if len(args) <= 0 {
		flag.Usage()
		return
	}

	var stdin io.Reader = os.Stdin
	var buf *bytes.Buffer
	var rw io.Reader
	if !terminal.IsTerminal(int(os.Stdin.Fd())) {
		buf = bytes.NewBuffer(make([]byte, 0, 1024))
		rw = ioutils.ReaderFunc(func(p []byte) (int, error) {
			n, err := os.Stdin.Read(p)
			if n > 0 {
				buf.Write(p[:n])
			}
			return n, err
		})
		stdin = rw
	}

	for i := 0; i < n; i++ {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdin = stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err == nil {
			return
		}

		if !silent {
			fmt.Fprintln(os.Stderr, err.Error()+", retry after "+strconv.Itoa(delay)+" seconds.")
		}

		time.Sleep(time.Duration(delay) * time.Second)

		if buf != nil {
			stdin = io.MultiReader(bytes.NewReader(buf.Bytes()), rw)
		}
	}

	os.Exit(1)
}
