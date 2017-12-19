%{
package main

import (
	"bufio"
	"errors"
	"fmt"
	"go/token"
	"strings"
	"unicode"

	"github.com/cznic/golex/lex"
	"github.com/mattn/komachi"
)

type lexer struct{
	*lex.Lexer
	empty   bool
	current byte
	err     error
	model   w2v.Model
	vector  *w2v.Vector
}

func clazz(r rune) int {
	if r < 0 || unicode.IsSpace(r) || r == '+' || r == '-' {
		return int(r)
	}
	return 0x80
}

func newLexer(model w2v.Model, src *bufio.Reader) *lexer {
	f := token.NewFileSet().AddFile("w2v", -1, 1<<31-1)
	lx, err := lex.New(f, bufio.NewReader(src), lex.RuneClass(clazz))
	if err != nil {
		panic(err)
	}
	return &lexer{Lexer: lx, model: model}
}

func (l *lexer) Error(e string) {
	l.err = errors.New(e)
}

func (l *lexer) Lex(lval *yySymType) int {
	c := l.Enter()
%}

%yyc c
%yyn c = l.Next()
%yym l.Mark()

wordChar      \x80
word          {wordChar}+

%%
	c = l.Rule0()

[ \t\r\n]+

{word}
	s := strings.TrimSpace(string(l.TokenBytes(nil)))
	if s != "" {
		if lval.value = l.model.Find(s); lval.value == nil {
			l.Error(fmt.Sprintf("%q not found", s))
			lval = nil
		}
		return VALUE
	}

%%
	if c, ok := l.Abort(); ok { return int(c) }
	goto yyAction
}
