package main

import (
	"bufio"
	"strings"

	"github.com/mattn/komachi"
)

type eval struct {
	model w2v.Model
}

func newEval(model w2v.Model) *eval {
	return &eval{
		model: model,
	}
}

func (e *eval) Do(s string) (*w2v.Vector, error) {
	lexer := newLexer(e.model, bufio.NewReader(strings.NewReader(s+"\n")))
	if r := yyParse(lexer); r == 0 {
		return lexer.vector, lexer.err
	}
	return nil, lexer.err
}
