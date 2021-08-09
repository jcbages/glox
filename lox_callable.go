package main

import "time"

type LoxCallable interface {
	Arity() int
	Call(intr *Interpreter, arguments []interface{}) (interface{}, error)
}

type FunctionLoxCallable struct {
	closure     *Environment
	declaration StmtFunction
}

type ClockLoxCallable struct{}

// Arity

func (lc FunctionLoxCallable) Arity() int {
	return len(lc.declaration.Parameters)
}

func (lc ClockLoxCallable) Arity() int {
	return 0
}

// Call

func (lc FunctionLoxCallable) Call(intr *Interpreter, arguments []interface{}) (interface{}, error) {
	environment := NewEnvironment(lc.closure)

	for i := 0; i < len(arguments); i += 1 {
		environment.Define(lc.declaration.Parameters[i].Lexeme, arguments[i])
	}

	err := intr.executeBlock(lc.declaration.Body, environment)
	if err == nil {
		return nil, nil
	} else if ret, ok := err.(Return); ok {
		return ret.Value, nil
	} else {
		return nil, err
	}
}

func (lc ClockLoxCallable) Call(intr *Interpreter, arguments []interface{}) (interface{}, error) {
	return float64(time.Now().Unix()), nil
}
