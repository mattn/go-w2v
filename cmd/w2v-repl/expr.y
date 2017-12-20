%{
package main

import (
	"github.com/mattn/go-w2v"
)

%}

%union{
	value *w2v.Vector
}

%token	VALUE

%left	'-' '+'

%type	<value>	VALUE, exp

%%
input:    /* empty */
        | input line
;

line:     '\n'
        | exp {
			yylex.(*lexer).vector = $1
		}
;

exp:      VALUE              { $$ = $1         }
        | exp '+' exp        { $$ = $1.Add($3) }
        | exp '-' exp        { $$ = $1.Sub($3) }
;
%%
