package annotation

import (
	"strings"
)

type Annotation struct {
	Name Name
	Body *Body
}

func (a Annotation) Start() Token {
	return a.Name.Token
}

type Name struct {
	Token Token
}

func (n Name) String() string {
	return n.Token.Value
}

type Body struct {
	Tokens []Token
}

func (v Body) Start() Token {
	return v.Tokens[0]
}

func (v Body) String() string {
	if len(v.Tokens) == 0 {
		return ""
	}

	var b strings.Builder
	for i, t := range v.Tokens {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(t.Value)
	}

	return b.String()
}

// Stencil returns the body padded with empty lines and leading spaces
// to match each token's line/column in the source.
//
//	# comment         →  ``
//	# @foo {          →  Pause, `       {`
//	#   "bar": 123    →  `    "bar": 123`
//	# }               →  `  }`
func (v Body) Stencil() string {
	var b strings.Builder
	line := 1
	for _, token := range v.Tokens {
		for line < token.Line {
			b.WriteString("\n")
			line++
		}
		b.WriteString(strings.Repeat(" ", token.Column-1))
		b.WriteString(token.Value)
	}
	return b.String()
}
