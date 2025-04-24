package main

import "fmt"

type LoxFunction struct {
	Declaration *FunctionStatement
	Closure     *Environment
	IsInitializer bool
}

func NewLoxFunction(declaration *FunctionStatement, closure *Environment, isInitializer bool) *LoxFunction {
	return &LoxFunction{
		Declaration: declaration,
		Closure:     closure,
		IsInitializer: isInitializer,
	}
}

func (l *LoxFunction) call(interpreter *Interpreter, arguments []interface{}) (result interface{}) {
	environment := NewEnclosedEnvironment(l.Closure)

	for index, param := range l.Declaration.Params {
		environment.define(param.Lexeme, arguments[index])
	}

	defer func() {
		if r := recover(); r != nil {
			if returnVal, ok := r.(*ReturnValue); ok {
                if l.IsInitializer {
                    result = l.Closure.getAt(0, "this")
                } else {
                    result = returnVal.Value
                }
			} else {
				panic(r)
			}
		}
	}()

	interpreter.executeBlock(l.Declaration.Body, environment)
	if (l.IsInitializer) {
		return l.Closure.getAt(0, "this")
	}
	return nil
}

func (l *LoxFunction) arity() int {
	return len(l.Declaration.Params)
}

func (l *LoxFunction) String() string {
	return fmt.Sprintf("<fn %s>", l.Declaration.Name.Lexeme)
}

func (l *LoxFunction) bind(instance *LoxInstance) *LoxFunction {
	environment := NewEnclosedEnvironment(l.Closure)
	environment.define("this", instance)
	return NewLoxFunction(l.Declaration, environment, l.IsInitializer)
}