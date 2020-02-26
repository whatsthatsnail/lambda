package environment

import (
	"lambda/ast"
)

type Environment struct {
	Values map[string]ast.Term
}

func NewEnvironment() Environment {
	env := Environment{}
	env.Values = make(map[string]ast.Term)
	return env
}

func (e *Environment) Get(name string) (interface{}, bool) {
	value, ok := e.Values[name]
	if ok {
		return value, ok
	}

	return nil, ok
}

func (e *Environment) Define(name string, term ast.Term) {
	e.Values[name] = term
}

func (e *Environment) Lookup(term ast.Term) (string, bool) {
	for k, v := range e.Values {
		if v == term {
			return k, true
		}
	}

	return "", false
}
