package main

import (
	// Uncomment this to pass the first stage
	"fmt"
	"io"
	"os"
	"strings"
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


func isWord(letter byte) bool {
    return isDigit(letter) || isLetter(letter) || letter == '_'
}

func isDigit(letter byte) bool {
    return letter >= '0' && letter <= '9'
}

func isLetter(letter byte) bool {
    return letter >= 'A' && letter <='Z' || letter >= 'a' && letter <= 'z'
}

func isSpecial(letter byte) bool {
    return strings.Contains("\\[]()^$", string(letter))
}

func isInGroup(group string) func (byte) bool {
    if strings.HasPrefix(group, "^") {
        return func (letter byte) bool {
            return !strings.Contains(group, string(letter))
        }
    } else {
        return func (letter byte) bool {
            return strings.Contains(group, string(letter))
        }
    }
}

type predicate func (byte) bool
func matchHere(line string, pattern string, here int) bool {
    p := 0

    var end bool
    if strings.HasSuffix(pattern, "$") {
        end = true
        pattern = pattern[:len(pattern)-1]
    }

    var functor predicate
    var inc int

    for l := here; l < len(line); {
        fmt.Printf("Line pointer now points at %v\n", l)
        fmt.Printf("Pattern pointer now points at %v\n", p)
        if pattern[p] == '\\' {
            if p + 1 >= len(pattern) {
                panic("pattern ended unexpectedly")
            }
            p += 1
            if pattern[p] == 'd' {
                // println("The pattern is \\d")
                functor = isDigit
            } else if pattern[p] == 'w' {
                // println("The pattern is \\w")
                functor = isWord
            }
            inc = 1
        } else if pattern[p] == '[' {
            closing := strings.Index(pattern[p:], string(']'))
            if closing == -1 {
                panic("no closing bracket")
            }
            group := pattern[p+1:closing+p] 
            fmt.Printf("Pattern is [%s]", group)
            // println(p+1, closing+p+1)
            functor = isInGroup(group)
            inc = closing + 2
        } else if pattern[p] == '+' {
            fmt.Printf("pattern is +")
            if p == len(pattern) - 1 {
                return true
            }
            for ;l < len(line); l++{
                if matchHere(line, pattern[p+1:], l) {
                    return true
                }
            }
            return false
        } else {
            fmt.Printf("No pattern, just letter")
            functor = isInGroup(string(pattern[p]))
            inc = 1
        }
        if !functor(line[l]) {
            return false
        } else {
            p += inc
        }
        l += 1
        if p >= len(pattern) {
            if end {
                return l == len(line)
            }
            return true
        }
    }
    return false
}

func matchLine(line []byte, pattern string) (bool, error) {
	if utf8.RuneCountInString(pattern) == 0 {
		return false, fmt.Errorf("unsupported pattern: %q", pattern)
	}


	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

    if strings.HasPrefix(pattern, "^") {
        return matchHere(string(line), pattern[1:], 0), nil
    }

    for i := range line {
        if matchHere(string(line), pattern, i) {
            print(true)
            return true, nil
        }
    }
    return false, nil
}
