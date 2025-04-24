package main

type LoxClass struct {
	Name    string
	Methods map[string]*LoxFunction
	Superclass *LoxClass
}

func NewLoxClass(name string, superclass *LoxClass, methods map[string]*LoxFunction) *LoxClass {
	return &LoxClass{
		Name:    name,
		Methods: methods,
		Superclass: superclass,
	}
}

func (l *LoxClass) arity() int {
	initializer := l.findMethod("init")
	if initializer == nil {
		return 0 
	}
	return initializer.arity()
}

func (l *LoxClass) call(interpreter *Interpreter, arguments []interface{}) interface{} {
	instance := NewLoxInstance(l)

	initializer := l.findMethod("init") 

	if initializer != nil {
		initializer.bind(instance).call(interpreter, arguments)
	}

	return instance
}

func (l *LoxClass) findMethod(name string) *LoxFunction {
    if method, ok := l.Methods[name]; ok {
        return method
    }

	if l.Superclass != nil { 
		return l.Superclass.findMethod(name)
	}
	
    return nil
}

func (l *LoxClass) String() string {
	return l.Name
}