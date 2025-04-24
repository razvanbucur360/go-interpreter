package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

type Interpreter struct {
	shouldPrintExpressions bool
	globals                *Environment
	environment            *Environment
	locals                 map[Expr]int
}

type ReturnValue struct {
	Value interface{}
}

func (i *Interpreter) interpret(statements []Stmt) {
	defer func() {
		if r := recover(); r != nil {
			if runtimeErr, ok := r.(RuntimeError); ok && !hasError {
				handleRuntimeError(runtimeErr)
			} else if !hasError {
				panic(r)
			}
		}
	}()

	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i *Interpreter) resolve(expr Expr, depth int) {
	i.locals[expr] = depth
}

func (i *Interpreter) execute(statement Stmt) {
	statement.Accept(i)
}

func (i *Interpreter) stringify(object interface{}) string {
	if object == nil {
		return "nil"
	}

	switch v := object.(type) {
	case float64:
		// Format the number to 2 decimal places first
		formatted := fmt.Sprintf("%.2f", v)
		// Remove trailing zeros and decimal point if there's no fractional part
		if strings.Contains(formatted, ".") {
			formatted = strings.TrimRight(formatted, "0")
			formatted = strings.TrimRight(formatted, ".")
		}
		return formatted
	default:
		return fmt.Sprintf("%v", object)
	}
}

// Expression visitor function implementations

func (i *Interpreter) visitBinaryExpr(expr *BinaryExpr) interface{} {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.TokenType {
	case MINUS:
		i.checkNumberOperands(left, expr.Operator, right)
		return left.(float64) - right.(float64)
	case STAR:
		i.checkNumberOperands(left, expr.Operator, right)
		return left.(float64) * right.(float64)
	case SLASH:
		i.checkNumberOperands(left, expr.Operator, right)
		return left.(float64) / right.(float64)
	case PLUS:
		if leftNum, leftOk := left.(float64); leftOk {
			if rightNum, rightOk := right.(float64); rightOk {
				return leftNum + rightNum
			}
		}
		if leftStr, leftOk := left.(string); leftOk {
			if rightStr, rightOk := right.(string); rightOk {
				return leftStr + rightStr
			}
		}
		panic(RuntimeError{Token: expr.Operator, Message: "Operands must be two numbers or two string"})
	case GREATER:
		i.checkNumberOperands(left, expr.Operator, right)
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		i.checkNumberOperands(left, expr.Operator, right)
		return left.(float64) >= right.(float64)
	case LESS:
		i.checkNumberOperands(left, expr.Operator, right)
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
		i.checkNumberOperands(left, expr.Operator, right)
		return left.(float64) <= right.(float64)
	case BANG_EQUAL:
		return !i.isEqual(left, right)
	case EQUAL_EQUAL:
		return i.isEqual(left, right)
	}

	return nil
}

func (i *Interpreter) visitUnaryExpr(expr *UnaryExpr) interface{} {
	right := i.evaluate(expr.Right)

	switch expr.Operator.TokenType {
	case MINUS:
		i.checkOperand(expr.Operator, right)
		return -right.(float64)
	case BANG:
		return !i.isTruthy(right)
	}

	return nil
}

func (i *Interpreter) visitGroupingExpr(expr *GroupingExpr) interface{} {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) visitLiteralExpr(expr *LiteralExpr) interface{} {
	return expr.Value
}

func (i *Interpreter) visitVariableExpr(expr *VariableExpr) interface{} {
	return i.lookUpVariable(expr.Name, expr)
}

func (i *Interpreter) visitAssignmentExpr(expr *AssignmentExpr) interface{} {
	value := i.evaluate(expr.Value)
	distance, exists := i.locals[expr]
	if exists {
		i.environment.assignAt(distance, expr.Name, value)
	} else {
		i.globals.assign(expr.Name, value)
	}
	return value
}

func (i *Interpreter) visitLogicalExpr(expr *LogicalExpr) interface{} {
	leftExpr := i.evaluate(expr.Left)
	if expr.Operator.TokenType == OR {
		if i.isTruthy(leftExpr) {
			return leftExpr
		}
	} else {
		if !i.isTruthy(leftExpr) {
			return leftExpr
		}
	}

	return i.evaluate(expr.Right)
}

func (i *Interpreter) visitCallExpr(expr *CallExpression) interface{} {
	callee := i.evaluate(expr.Callee)

	arguments := make([]interface{}, 0, len(expr.Arguments))
	for _, arg := range expr.Arguments {
		arguments = append(arguments, i.evaluate(arg))
	}

	function, ok := callee.(LoxCallable)
	if !ok {
		panic(RuntimeError{
			Token:   expr.Parenthesis,
			Message: "Can only call functions and classes",
		})
	}

	if len(arguments) != function.arity() {
		panic(RuntimeError{
			Token:   expr.Parenthesis, // Make sure this matches your struct field name
			Message: fmt.Sprintf("Expected %d arguments but got %d.", function.arity(), len(arguments)),
		})
	}

	return function.call(i, arguments)
}

func (i *Interpreter) visitGetExpr(expr *GetExpression) interface{} {
	object := i.evaluate(expr.Object)

	if instance, ok := object.(*LoxInstance); ok {
		return instance.get(expr.Name)
	}

	panic(RuntimeError{
		Token: expr.Name,
		Message: "Only instances have properties.",
	})
}

func (i *Interpreter) visitSetExpr(expr *SetExpression) interface{} {
	object := i.evaluate(expr.Object)

	if _, ok := object.(*LoxInstance); !ok {
		panic(RuntimeError{
			Token: expr.Name,
			Message: "Only instances have fields.",
		})
	}
	
	value := i.evaluate(expr.Value)
	object.(*LoxInstance).set(expr.Name, value)
	return value 
}

func (i *Interpreter) visitThisExpr(expr *ThisExpr) interface{} {
	return i.lookUpVariable(expr.Keyword, expr)
}

func (i *Interpreter) visitSuperExpr(expr *SuperExpr) interface{} {
	distance := i.locals[expr]
	superclass := i.environment.getAt(distance, "super").(*LoxClass)
	object := i.environment.getAt(distance - 1, "this").(*LoxInstance)
	method := superclass.findMethod(expr.Method.Lexeme) 

	if method == nil {
		panic(&RuntimeError{
			Token: expr.Method,
			Message: "Undefined property '" + expr.Method.Lexeme + "'.",
		})
	}

	return method.bind(object)
}

// ----------------------------------------------

// Statement visitor function implementations

func (i *Interpreter) visitPrintStmt(stmt *PrintStatement) interface{} {
	value := i.evaluate(stmt.Value)
	fmt.Println(i.stringify(value))
	return nil
}

func (i *Interpreter) visitExpressionStmt(stmt *ExpressionStatement) interface{} {
	value := i.evaluate(stmt.Expression)
	if i.shouldPrintExpressions {
		fmt.Println(i.stringify(value))
	}
	return nil
}

func (i *Interpreter) visitVarStmt(stmt *VarStatement) interface{} {
	var value interface{}

	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}

	i.environment.define(stmt.Name.Lexeme, value)

	return nil
}

func (i *Interpreter) visitBlockStmt(stmt *Block) interface{} {
	enclosingEnv := i.environment
	blockEnv := NewEnclosedEnvironment(enclosingEnv)
	i.executeBlock(stmt.Statements, blockEnv)
	return nil
}

func (i *Interpreter) visitIfStmt(stmt *IfStatement) interface{} {
	if i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		i.execute(stmt.ElseBranch)
	}
	return nil
}

func (i *Interpreter) visitWhileStmt(stmt *WhileStatement) interface{} {
	for i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.Body)
	}
	return nil
}

func (i *Interpreter) visitFunctionStmt(stmt *FunctionStatement) interface{} {
	function := NewLoxFunction(stmt, i.environment, false)
	i.environment.define(stmt.Name.Lexeme, function)
	return nil
}

func (i *Interpreter) visitReturnStmt(stmt *ReturnStatement) interface{} {
	var value interface{} = nil

	if stmt.Value != nil {
		value = i.evaluate(stmt.Value)
	}

	panic(&ReturnValue{
		Value: value, 
	})
}

func (i *Interpreter) visitClassStmt(stmt *ClassStatement) interface{} {
    var superclass *LoxClass
    if stmt.Superclass != nil {
        sc := i.evaluate(stmt.Superclass)
        var ok bool
        superclass, ok = sc.(*LoxClass)
        if !ok {
            panic(RuntimeError{ 
                Token:   stmt.Superclass.Name,
                Message: "Superclass must be a class",
            })
        }
    }

    i.environment.define(stmt.Name.Lexeme, nil)

    enclosingEnv := i.environment
    if stmt.Superclass != nil {
        i.environment = NewEnclosedEnvironment(i.environment)
        i.environment.define("super", superclass)
    }

    methods := make(map[string]*LoxFunction)
    for _, method := range stmt.Methods {
        function := NewLoxFunction(method, i.environment, method.Name.Lexeme == "init")
        methods[method.Name.Lexeme] = function
    }

    klass := NewLoxClass(stmt.Name.Lexeme, superclass, methods)

    if stmt.Superclass != nil {
        i.environment = enclosingEnv
    }

    i.environment.assign(stmt.Name, klass)
    return nil
}

// ----------------------------------------------

func (i *Interpreter) executeBlock(statements []Stmt, environment *Environment) {
	previousEnvironment := i.environment
	i.environment = environment

	defer func() {
		i.environment = previousEnvironment
	}()

	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i *Interpreter) checkNumberOperands(leftOperand interface{}, operator Token, rightOperand interface{}) {
	_, leftOk := leftOperand.(float64)
	_, rightOk := rightOperand.(float64)

	if leftOk && rightOk {
		return
	}

	panic(RuntimeError{Token: operator, Message: "Operands must be a number."})

}

func (i *Interpreter) checkOperand(operator Token, operand interface{}) {
	_, ok := operand.(float64)
	if ok {
		return
	}
	panic(RuntimeError{Token: operator, Message: "Operand must be a number."})
}

func (i *Interpreter) evaluate(expr Expr) interface{} {
	return expr.Accept(i)
}

func (i *Interpreter) isTruthy(object interface{}) bool {
	if object == nil {
		return false
	}

	if b, ok := object.(bool); ok {
		return b
	}

	return true
}

func (i *Interpreter) isEqual(left interface{}, right interface{}) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil {
		return false
	}

	return reflect.DeepEqual(left, right)
}

func (i *Interpreter) lookUpVariable(name Token, expr Expr) interface{} {
	distance, exists := i.locals[expr]
	if exists {
		return i.environment.getAt(distance, name.Lexeme)
	} else {
		return i.globals.get(name)
	}
}

type RuntimeError struct {
	Token   Token
	Message string
}

// Implement the error interface for RuntimeError
func (e RuntimeError) Error() string {
	return e.Message + " at " + e.Token.Lexeme
}

func handleRuntimeError(err RuntimeError) {
	fmt.Fprintf(os.Stderr, "%s\n[line %d]", err.Message, err.Token.Line)
	hasRuntimeError = true
}
