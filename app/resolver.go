package main

import "fmt"

type Resolver struct {
	Interpreter     *Interpreter
	Scopes          []map[string]bool // string for var/func name and bool for wether it's been defined or not. Initially we only declare and only after a safe check we define
	CurrentFunction FunctionType
	CurrentClass ClassType
}

func NewResolver(intepreter *Interpreter) *Resolver {
	r :=  &Resolver{
		Interpreter:     intepreter,
		Scopes:          make([]map[string]bool, 0),
		CurrentFunction: FUNCTION_NONE,
		CurrentClass: CLASS_NONE,
	}
	r.beginScope()
	return r
}

func (r *Resolver) beginScope() {
	r.Scopes = append(r.Scopes, make(map[string]bool))
}

func (r *Resolver) resolve(statements []Stmt) {
	for _, statement := range statements {
		r.resolveStatement(statement)
	}
}

func (r *Resolver) endScope() {
	r.Scopes = r.Scopes[:len(r.Scopes)-1]
}

func (r *Resolver) resolveStatement(stmt Stmt) {
    if stmt == nil {
        return
    }
    stmt.Accept(r)
}

func (r *Resolver) resolveExpression(expr Expr) {
	expr.Accept(r)
}

func (r *Resolver) visitBlockStmt(stmt *Block) interface{} {
    r.beginScope()
    r.resolve(stmt.Statements)
    r.endScope()
    return nil
}

func (r *Resolver) visitVarStmt(stmt *VarStatement) interface{} {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolveExpression(stmt.Initializer)
	}
	r.define(stmt.Name)
	return nil
}

func (r *Resolver) declare(name Token) {
    if len(r.Scopes) == 0 {
        return
    }
    scope := r.Scopes[len(r.Scopes)-1]
	if len(r.Scopes) > 1{
		if _, exists := scope[name.Lexeme]; exists {
			r.error(name, "Already a variable with this name in this scope.")
		}
	}

    if len(r.Scopes) > 1 {
        scope[name.Lexeme] = false
    }
}

func (r *Resolver) define(name Token) {
    if len(r.Scopes) == 0 {
        return
    }
    r.Scopes[len(r.Scopes)-1][name.Lexeme] = true
}

func (r *Resolver) visitVariableExpr(expr *VariableExpr) interface{} {
    if len(r.Scopes) > 1 {
        if initialized, exists := r.Scopes[len(r.Scopes)-1][expr.Name.Lexeme]; exists && !initialized {
			r.error(expr.Name, "Can't read local variable in its own initializer.")
        }
    }
    r.resolveLocal(expr, expr.Name)
    return nil
}

func (r *Resolver) resolveLocal(expr Expr, name Token) {
	for i := len(r.Scopes) - 1; i >= 0; i-- {
		if _, exists := r.Scopes[i][name.Lexeme]; exists {
			r.Interpreter.resolve(expr, len(r.Scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) visitAssignmentExpr(expr *AssignmentExpr) interface{} {
	r.resolveExpression(expr.Value)
	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) visitGetExpr(expr *GetExpression) interface{} {
	r.resolveExpression(expr.Object)
	return nil 
}

func (r *Resolver) visitSetExpr(expr *SetExpression) interface{} {
	r.resolveExpression(expr.Object)
	r.resolveExpression(expr.Value)
	return nil 
}

func (r *Resolver) visitFunctionStmt(stmt *FunctionStatement) interface{} {
	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.resolveFunction(stmt, FUNCTION_FUNCTION)
	return nil
}

func (r *Resolver) resolveFunction(function *FunctionStatement, functionType FunctionType) {
    enclosingFunction := r.CurrentFunction
    r.CurrentFunction = functionType

    r.beginScope()
    for _, param := range function.Params {
        r.declare(param)
        r.define(param)
    }
    r.resolve(function.Body)
    r.endScope()
    
    r.CurrentFunction = enclosingFunction
}

func (r *Resolver) visitExpressionStmt(stmt *ExpressionStatement) interface{} {
    if stmt.Expression != nil {
        r.resolveExpression(stmt.Expression)
    }
    return nil
}

func (r *Resolver) visitIfStmt(stmt *IfStatement) interface{} {
	r.resolveExpression(stmt.Condition)
	r.resolveStatement(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolveStatement(stmt.ElseBranch)
	}
	return nil
}

func (r *Resolver) visitPrintStmt(stmt *PrintStatement) interface{} {
	r.resolveExpression(stmt.Value)
	return nil
}

func (r *Resolver) visitReturnStmt(stmt *ReturnStatement) interface{} {
	if r.CurrentFunction == FUNCTION_NONE {
		r.error(stmt.Keyword, "Can't return from top-level code.")
	}
	if stmt.Value != nil {
		if r.CurrentFunction == FUNCTION_INITIALIZER {
			r.error(stmt.Keyword, "Can't return a value from an initializer.")
		}
		r.resolveExpression(stmt.Value)
	}
	return nil
}

func (r *Resolver) visitWhileStmt(stmt *WhileStatement) interface{} {
	r.resolveExpression(stmt.Condition)
	r.resolveStatement(stmt.Body)
	return nil
}

func (r *Resolver) visitClassStmt(stmt *ClassStatement) interface{} {
	enclosingClass := r.CurrentClass
	r.CurrentClass = CLASS_CLASS

	r.declare(stmt.Name)
	r.define(stmt.Name)

	if stmt.Superclass != nil {
		if stmt.Name == stmt.Superclass.Name {
			r.error(stmt.Superclass.Name, "A class can't inherit from itself.")
		}
	}

	if stmt.Superclass != nil {
		r.CurrentClass = CLASS_SUBCLASS
		r.resolveExpression(stmt.Superclass)
	}

	if stmt.Superclass != nil {
		r.beginScope()
		r.Scopes[len(r.Scopes) -1]["super"] = true 
	}

	r.beginScope()
	r.Scopes[len(r.Scopes) - 1]["this"] = true   
	
	for _, method := range stmt.Methods {
		declaration := FUNCTION_METHOD 
		if method.Name.Lexeme == "init" {
			declaration = FUNCTION_INITIALIZER
		}
		r.resolveFunction(method, declaration)
	}

	r.endScope()
	if stmt.Superclass != nil { 
		r.endScope()
	}

	r.CurrentClass = enclosingClass
	return nil 
}

func (r *Resolver) visitBinaryExpr(expr *BinaryExpr) interface{} {
	r.resolveExpression(expr.Left)
	r.resolveExpression(expr.Right)
	return nil
}

func (r *Resolver) visitCallExpr(expr *CallExpression) interface{} {
	r.resolveExpression(expr.Callee)
	for _, expression := range expr.Arguments {
		r.resolveExpression(expression)
	}
	return nil
}

func (r *Resolver) visitGroupingExpr(expr *GroupingExpr) interface{} {
	r.resolveExpression(expr.Expression)
	return nil
}

func (r *Resolver) visitLiteralExpr(expr *LiteralExpr) interface{} {
	return nil
}

func (r *Resolver) visitLogicalExpr(expr *LogicalExpr) interface{} {
	r.resolveExpression(expr.Left)
	r.resolveExpression(expr.Right)
	return nil
}

func (r *Resolver) visitUnaryExpr(expr *UnaryExpr) interface{} {
	r.resolveExpression(expr.Right)
	return nil
}

func (r *Resolver) visitThisExpr(expr *ThisExpr) interface{} {
	if r.CurrentClass == CLASS_NONE {
		r.error(expr.Keyword, "Can't use 'this' outside of a class.")
	}
	r.resolveLocal(expr, expr.Keyword)
	return nil 
}

func (r *Resolver) visitSuperExpr(expr *SuperExpr) interface{} {
	if r.CurrentClass == CLASS_NONE {
		r.error(expr.Keyword, "Can't use 'super' outside of a class.")
	} else if r.CurrentClass == CLASS_CLASS {
		r.error(expr.Keyword, "Can't use 'super' in a class with no superclass.")
	} 
	r.resolveLocal(expr, expr.Keyword)
	return nil 
}

type FunctionType int

const (
	FUNCTION_NONE FunctionType = iota
	FUNCTION_FUNCTION
	FUNCTION_INITIALIZER 
	FUNCTION_METHOD
)

type ClassType int 

const (
	CLASS_NONE ClassType = iota 
	CLASS_CLASS 
	CLASS_SUBCLASS
)

func (r *Resolver) error(token Token, message string) {
    report(token.Line, fmt.Sprintf(" at '%s'", token.Lexeme), message)
}