package main

import (
	"errors"
	"strconv"
	"strings"
)

type symbolType uintptr

const (
	// symbolTypeConst fields
	//
	//   0: const eval
	symbolTypeConst symbolType = iota
	// symbolTypeFunc fields
	//
	//   0: _para1_name, _para2_name, ...
	//   1: replaceable eval
	symbolTypeFunc
)

type symbol struct {
	stype  symbolType
	name   string
	fields []string
}

func (s *symbol) replace(txt string, args ...string) (string, error) {
	switch s.stype {
	case symbolTypeConst:
		return strings.ReplaceAll(txt, s.name, s.fields[0]), nil
	case symbolTypeFunc:
		//TODO: finish
		return "", nil
	}
	return "", errors.New("unsupported symbol type " + strconv.Itoa(int(s.stype)))
}
