package main

import "fmt"

type LoxInstance struct {
	Klass *LoxClass
	Fields map[string]interface{}
}

func NewLoxInstance(klass *LoxClass) *LoxInstance {
	return &LoxInstance{
		Klass: klass,
		Fields: make(map[string]interface{}),
	}
}

func (l *LoxInstance) get(name Token) interface{} {
	if _, exists := l.Fields[name.Lexeme]; exists {
		return l.Fields[name.Lexeme]
	}

	method := l.Klass.findMethod(name.Lexeme) 
	
	if method != nil {
		return method.bind(l)
	}

	panic(RuntimeError{
		Token: name,
		Message: fmt.Sprintf("Undefined property %s.", name.Lexeme),
	})
}

func (l *LoxInstance) set(name Token, value interface{}) {
	l.Fields[name.Lexeme] = value 
}

func (l *LoxInstance) String() string {
	return fmt.Sprintf("%s instance", l.Klass.Name)
}


