package main

import "time"

type LoxClock struct {}

func(l *LoxClock) arity() int {
	return 0
}

func(l *LoxClock) call(interpreter *Interpreter, arguments []interface{}) interface{} {
	return float64(time.Now().UnixNano()) / 1e9
}

func (l *LoxClock) String() string {
    return "<native fn>"
}

func NewInterpreter(shouldPrintExpressions bool) *Interpreter {
	globals := NewEnvironment()
	globals.define("clock", &LoxClock{})
	return &Interpreter{
		shouldPrintExpressions: shouldPrintExpressions,
		globals: globals,
		environment: globals,
		locals: make(map[Expr]int),
	}
}