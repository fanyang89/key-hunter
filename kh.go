package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/schollz/progressbar/v3"
)

var flagBlockSize = flag.Int("block-size", 4096, "logical block size")
var flagInputFile = flag.String("i", "-", "input")

func doMain() error {
	flag.Parse()
	if *flagInputFile == "" {
		flag.PrintDefaults()
		return nil
	}

	r, err := createReader(*flagInputFile)
	if err != nil {
		return err
	}

	bar := progressbar.DefaultBytes(-1, "Scanning")
	buffer := make([]byte, *flagBlockSize)

	for {
		var n int
		n, err = r.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		s := string(buffer[:n])
		if isKeyFile(s) {
			var tf *os.File
			tf, err = os.CreateTemp(".", "key-hunter-")
			if err != nil {
				return err
			}
			bar.Describe(fmt.Sprintf("Write to %s", tf.Name()))
			_, _ = tf.WriteString(s)
			_ = tf.Close()
		}
		_ = bar.Add(n)
	}

	return bar.Finish()
}

func isKeyFile(s string) bool {
	var fst string
	p := strings.IndexRune(s, '\n')
	if p == -1 {
		fst = s
	} else {
		fst = s[:p]
	}
	if strings.HasPrefix(fst, "ssh-") {
		return true
	}
	return strings.Contains(fst, "BEGIN OPENSSH PRIVATE KEY") &&
		strings.Contains(s, "END OPENSSH PRIVATE KEY")
}

func main() {
	err := doMain()
	if err != nil {
		panic(err)
	}
}
