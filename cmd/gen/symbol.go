package main

import (
	"bytes"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	errIsConstReplace = errors.New("is const replace")
	errNoSuchSymbol   = errors.New("no such sybmol")
)

var symtab = symbolTable{
	"ZE_APICALL":   &symbol{symbolTypeConst, symbolSubTypeDefine, "ZE_APICALL", []string{""}},
	"ZE_APIEXPORT": &symbol{symbolTypeConst, symbolSubTypeDefine, "ZE_APIEXPORT", []string{""}},
	"ZE_DLLEXPORT": &symbol{symbolTypeConst, symbolSubTypeDefine, "ZE_DLLEXPORT", []string{""}},
	"~(0ULL":       &symbol{symbolTypeConst, symbolSubTypeDefine, "~(0ULL", []string{"(^uint64(0)"}},
}

type symbolTable map[string]*symbol

func (st symbolTable) apply(t string) string {
	for _, s := range st {
		t = s.replace(t)
	}
	return t
}

func (st symbolTable) contains(name string) bool {
	_, ok := st[name]
	return ok
}

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

type symbolSubType uintptr

const (
	symbolSubTypeDefine symbolSubType = iota
	symbolSubTypeEmptyStruct
	symbolSubTypeLargeStruct
	symbolSubTypeEnum
	symbolSubTypeFuncPtr
)

type symbol struct {
	stype  symbolType
	sstype symbolSubType
	name   string
	fields []string
}

func newSymbolConst(name, val string, sstype symbolSubType) *symbol {
	return &symbol{
		stype:  symbolTypeConst,
		sstype: sstype,
		name:   name,
		fields: []string{val},
	}
}

func newSymbolFunc(name, paras, evals string, sstype symbolSubType) *symbol {
	return &symbol{
		stype:  symbolTypeFunc,
		sstype: sstype,
		name:   name,
		fields: []string{paras, evals},
	}
}

func (s *symbol) extract1stFunc(txt string) (args []string, a, b int, err error) {
	if s.stype == symbolTypeConst {
		return nil, 0, 0, errIsConstReplace
	}
	a = strings.Index(txt, s.name)
	if a < 0 {
		return nil, 0, 0, errNoSuchSymbol
	}
	str, off, err := getInsideRoundBrakets(txt[a:])
	if err != nil {
		return nil, 0, 0, err
	}
	args = strings.Split(str, ",")
	for i, arg := range args {
		args[i] = strings.TrimSpace(arg)
	}
	return args, a, a + off, nil
}

func (s *symbol) replace(txt string) string {
	switch s.stype {
	case symbolTypeConst:
		escapes := regexp.QuoteMeta(s.name)
		re := regexp.MustCompile(`(` + escapes + `[^\w_])|(` + escapes + `$)`)
		return string(re.ReplaceAllFunc([]byte(txt), func(b []byte) []byte {
			last := rune(b[len(b)-1])
			if unicode.IsLetter(last) || last == '_' {
				return []byte(s.fields[0])
			}
			buf := bytes.NewBuffer(make([]byte, 0, 128))
			buf.WriteString(s.fields[0])
			buf.WriteByte(b[len(b)-1])
			return buf.Bytes()
		}))
	case symbolTypeFunc:
		paras := strings.Split(s.fields[0], ",")
		txts := []string{}
		for {
			args, a, b, err := s.extract1stFunc(txt)
			if err == errNoSuchSymbol {
				txts = append(txts, txt)
				return strings.Join(txts, "")
			}
			if len(paras) != len(args) {
				panic("args " + strings.Join(args, ", ") + " count " + strconv.Itoa(len(args)) + " is different from recorded " + s.fields[0])
			}
			n := len(txts)
			txts = append(txts, []string{txt[:a], "/* ", txt[a:b], ") */(", s.fields[1]}...)
			for i, p := range paras {
				txts[n+4] = strings.ReplaceAll(txts[n+4], strings.TrimSpace(p), args[i])
			}
			txt = txt[b:]
		}
	}
	panic("unsupported symbol type " + strconv.Itoa(int(s.stype)))
}
