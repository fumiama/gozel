package main

import (
	"bufio"
	"errors"
	"strings"
)

func skip2endif(scan *bufio.Scanner, ln int) int {
	depth := 1
	for scan.Scan() {
		ln++
		t := scan.Text()
		switch {
		case strings.HasPrefix(t, "#endif"):
			depth--
		case strings.HasPrefix(t, "#if"):
			depth++
		default:
		}
		if depth <= 0 {
			break
		}
	}
	return ln
}

func getinside0brakets(txt string) (string, int, error) {
	depth := 0
	a := 0
	for i, t := range txt {
		switch t {
		case '(':
			if depth == 0 {
				a = i + 1
			}
			depth++
		case ')':
			depth--
			if depth <= 0 {
				return txt[a:i], i, nil
			}
		}
	}
	return "", 0, errors.New("no round brakets pair in " + txt)
}
