package main

import "fmt"

type Environment struct {
	Enclosing *Environment
	Values    map[string]interface{}
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		Enclosing: enclosing,
		Values:    make(map[string]interface{}),
	}
}

func (env *Environment) Define(name string, value interface{}) {
	env.Values[name] = value
}

func (env *Environment) Get(token Token) (interface{}, error) {
	if value, ok := env.Values[token.Lexeme]; ok {
		return value, nil
	} else if env.Enclosing != nil {
		return env.Enclosing.Get(token)
	} else {
		return nil, RuntimeError{
			Token:   token,
			Message: fmt.Sprintf("Undefined variable '%v'", token.Lexeme),
		}
	}
}

func (env *Environment) Assign(token Token, value interface{}) error {
	if _, ok := env.Values[token.Lexeme]; ok {
		env.Values[token.Lexeme] = value
		return nil
	} else if env.Enclosing != nil {
		return env.Assign(token, value)
	} else {
		return RuntimeError{
			Token:   token,
			Message: fmt.Sprintf("Undefined variable '%v'", token.Lexeme),
		}
	}
}
