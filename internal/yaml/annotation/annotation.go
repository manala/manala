package annotation

import (
	"strings"
)

type Annotation struct {
	Name  Name
	Value Value
}

type Name struct {
	Token Token
}

func (n Name) String() string {
	return n.Token.Value
}

type Value struct {
	Tokens []Token
}

func (v Value) String() string {
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
