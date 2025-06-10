//go:build std

package main

import (
	"io"
	"os"
)

func createReader(input string) (io.Reader, error) {
	if input == "-" {
		return os.Stdin, nil
	}
	return os.Open(*flagInputFile)
}
