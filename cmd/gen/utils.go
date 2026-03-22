package main

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
	"unicode"
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

func getInsideRoundBrakets(txt string) (string, int, error) {
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

func get1sentence(firstln string, scan *bufio.Scanner, ln int) (string, int) {
	if strings.Contains(firstln, ";") && !strings.HasPrefix(strings.TrimSpace(firstln), "//") {
		return firstln, ln
	}
	bracedepth := 0
	sb := strings.Builder{}
	sb.WriteString(firstln)
	for scan.Scan() {
		sb.WriteString("\n")
		t := scan.Text()
		ln++
		if strings.Contains(t, "{") {
			bracedepth++
		}
		if strings.Contains(t, "}") {
			bracedepth--
		}
		sb.WriteString(t)
		content, _, _ := strings.Cut(t, "//")
		if strings.Contains(content, ";") && bracedepth == 0 {
			return sb.String(), ln
		}
	}
	return "", -1
}

func scanln(name string, scan *bufio.Scanner, ln *int) (s string, isfin bool) {
	if !scan.Scan() {
		panic(fmt.Sprintf("%s L%d: unexpected EOF", name, *ln))
	}
	(*ln)++
	s = scan.Text()
	content, _, _ := strings.Cut(s, "//")
	isfin = strings.Contains(content, ";")
	return
}

func us2camel(t string) string {
	sb := strings.Builder{}
	for s := range strings.SplitSeq(t, "_") {
		rs := []rune(s)
		sb.WriteRune(unicode.ToUpper(rs[0]))
		for _, r := range rs[1:] {
			sb.WriteRune(unicode.ToLower(r))
		}
	}
	return sb.String()
}

func trimEmptyStringArray(arr []string) []string {
	news := make([]string, 0, len(arr))
	for _, s := range arr {
		if s == "" {
			continue
		}
		news = append(news, s)
	}
	return news
}
