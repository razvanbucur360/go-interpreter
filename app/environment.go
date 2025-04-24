package main

import "fmt"

type Environment struct {
	values    map[string]interface{}
	enclosing *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		values:    make(map[string]interface{}),
		enclosing: nil,
	}
}

func NewEnclosedEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		values:    make(map[string]interface{}),
		enclosing: enclosing,
	}
}

func (e *Environment) define(name string, value interface{}) {
	if e.values == nil {
		e.values = make(map[string]interface{})
	}
	e.values[name] = value
}

func (e *Environment) get(name Token) interface{} {

	if value, ok := e.values[name.Lexeme]; ok {
		return value
	}

	if e.enclosing != nil {
		return e.enclosing.get(name)
	}

	panic(RuntimeError{
		Token:   name,
		Message: fmt.Sprintf("Undefined variable '%s'", name.Lexeme),
	})
}

func (e *Environment) getAt(distance int, name string) interface{} {
	return e.ancestor(distance).values[name]
}

func (e *Environment) ancestor(distance int) *Environment {
	environment := e
	for i := 0; i < distance; i++ {
		environment = environment.enclosing
	}
	return environment
}

func (e *Environment) assign(name Token, value interface{}) {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return
	}

	if e.enclosing != nil {
		e.enclosing.assign(name, value)
		return
	}

	panic(RuntimeError{
		Token:   name,
		Message: "Undefined variable '" + name.Lexeme + "'.",
	})
}

func (e *Environment) assignAt(distance int, name Token, value interface{}) {
	e.ancestor(distance).values[name.Lexeme] = value
}
