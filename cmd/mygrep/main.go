package main

import (
	// Uncomment this to pass the first stage
	"bytes"
	"fmt"
	"io"
	"os"
	"unicode/utf8"
)

// Usage: echo <input_text> | your_program.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin) // assume we're only dealing with a single line
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}

	ok, err := matchLine(line, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}

	// default exit code is 0 which means success
}

func matchWord(letter rune) bool {
    if !(letter >= 'A' && letter <='Z' || letter >= 'a' && letter <= 'z' || letter >= '0' && letter <= '9' || letter == '_') {
        return false
    }
    return true
}

func matchLine(line []byte, pattern string) (bool, error) {
	if utf8.RuneCountInString(pattern) == 0 {
		return false, fmt.Errorf("unsupported pattern: %q", pattern)
	}

	var ok bool

	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

    if pattern == "\\d" {
        ok = bytes.ContainsAny(line, "0123456789")
    } else if pattern == "\\w" {
        ok = bytes.ContainsFunc(line, matchWord)
    } else if pattern[0] == '[' && pattern[len(pattern)-1] == ']' {
        ok = bytes.ContainsAny(line, pattern[1:len(pattern)-1])
        if len(pattern) >= 3 && pattern[1] == '^' {
            ok = !ok
        } 
    } else {
        ok = bytes.ContainsAny(line, pattern)
    } 
    println(ok)

	return ok, nil
}
