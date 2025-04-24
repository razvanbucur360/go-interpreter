package main

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (a *AstPrinter) Print(statements []Stmt) string {
    var sb strings.Builder
    for _, stmt := range statements {
        printed := a.printStmt(stmt)
        sb.WriteString(printed)
    }
    return sb.String()
}

func (a *AstPrinter) printStmt(stmt Stmt) string {
    if stmt == nil {
        return "nil"
    }
    
    switch s := stmt.(type) {
    case *Block:
        return a.visitBlockStmt(s).(string)
    case *VarStatement:
        return a.visitVarStmt(s).(string)
    case *FunctionStatement:
        return a.visitFunctionStmt(s).(string)
    case *ExpressionStatement:
        printed := a.visitExpressionStmt(s).(string)
        return printed
    case *IfStatement:
        return a.visitIfStmt(s).(string)
    case *PrintStatement:
        return a.visitPrintStmt(s).(string)
    case *ReturnStatement:
        return a.visitReturnStmt(s).(string)
    case *WhileStatement:
        return a.visitWhileStmt(s).(string)
    default:
        return fmt.Sprintf("(unknown %T)", stmt)
    }
}

// Updated visitBlockStmt to handle nested statements better
func (a *AstPrinter) visitBlockStmt(stmt *Block) interface{} {
    var sb strings.Builder
    sb.WriteString("(block")
    for _, s := range stmt.Statements {
        printed := a.printStmt(s)
        if printed != "nil" {
            sb.WriteString(printed)
        }
    }
    sb.WriteString(")")
    return sb.String()
}

// Ensure all visit methods return strings
func (a *AstPrinter) visitVarStmt(stmt *VarStatement) interface{} {
    if stmt.Initializer == nil {
        return fmt.Sprintf("(var %s)", stmt.Name.Lexeme)
    }
    return fmt.Sprintf("(var %s %s)", stmt.Name.Lexeme, a.printExpr(stmt.Initializer))
}

func (a *AstPrinter) visitFunctionStmt(stmt *FunctionStatement) interface{} {
    var sb strings.Builder
    sb.WriteString("(fun ")
    sb.WriteString(stmt.Name.Lexeme)
    sb.WriteString(" (")
    
    for i, param := range stmt.Params {
        if i > 0 {
            sb.WriteString(" ")
        }
        sb.WriteString(param.Lexeme)
    }
    sb.WriteString(")")
    
    // Print function body
    body := a.printStmt(&Block{Statements: stmt.Body})
    sb.WriteString(" ")
    sb.WriteString(body)
    sb.WriteString(")")
    
    return sb.String()
}

func (a *AstPrinter) visitExpressionStmt(stmt *ExpressionStatement) interface{} {
    if stmt.Expression == nil {
        return "nil"
    }
    return a.printExpr(stmt.Expression)
}

func (a *AstPrinter) visitIfStmt(stmt *IfStatement) interface{} {
	if stmt.ElseBranch == nil {
		return fmt.Sprintf("(if %s %s)",
			a.printExpr(stmt.Condition),
			a.printStmt(stmt.ThenBranch))
	}
	return fmt.Sprintf("(if %s %s %s)",
		a.printExpr(stmt.Condition),
		a.printStmt(stmt.ThenBranch),
		a.printStmt(stmt.ElseBranch))
}

func (a *AstPrinter) visitPrintStmt(stmt *PrintStatement) interface{} {
	return fmt.Sprintf("(print %s)", a.printExpr(stmt.Value))
}

func (a *AstPrinter) visitReturnStmt(stmt *ReturnStatement) interface{} {
	if stmt.Value == nil {
		return "(return)"
	}
	return fmt.Sprintf("(return %s)", a.printExpr(stmt.Value))
}

func (a *AstPrinter) visitWhileStmt(stmt *WhileStatement) interface{} {
	return fmt.Sprintf("(while %s %s)",
		a.printExpr(stmt.Condition),
		a.printStmt(stmt.Body))
}

// Expression visitors
func (a *AstPrinter) visitBinaryExpr(expr *BinaryExpr) interface{} {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *AstPrinter) visitGroupingExpr(expr *GroupingExpr) interface{} {
	return a.parenthesize("group", expr.Expression)
}

func (a *AstPrinter) visitLiteralExpr(expr *LiteralExpr) interface{} {
    if expr.Value == nil {
        return "nil"
    }
    switch v := expr.Value.(type) {
    case float64:
        if v == float64(int(v)) {
            return fmt.Sprintf("%.1f", v)
        }
        return fmt.Sprintf("%g", v)
    case string:
        return v // Remove strconv.Quote here
    case bool:
        if v {
            return "true"
        }
        return "false"
    default:
        return fmt.Sprintf("%v", expr.Value)
    }
}

func (a *AstPrinter) visitUnaryExpr(expr *UnaryExpr) interface{} {
	return a.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (a *AstPrinter) visitVariableExpr(expr *VariableExpr) interface{} {
	return expr.Name.Lexeme
}

func (a *AstPrinter) visitAssignmentExpr(expr *AssignmentExpr) interface{} {
	return fmt.Sprintf("(assign %s %s)", expr.Name.Lexeme, a.printExpr(expr.Value))
}

func (a *AstPrinter) visitLogicalExpr(expr *LogicalExpr) interface{} {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *AstPrinter) visitCallExpr(expr *CallExpression) interface{} {
	var sb strings.Builder
	sb.WriteString(a.printExpr(expr.Callee))
	sb.WriteString("(")
	for i, arg := range expr.Arguments {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(a.printExpr(arg))
	}
	sb.WriteString(")")
	return sb.String()
}

//Actual implementation not needed unless printing of the expressions is required. The function declarations are needed in order fot the AstPrinter to be an ExprVisitor
func (a *AstPrinter) visitGetExpr(expr *GetExpression) interface{} {
    return nil
}

func (a *AstPrinter) visitSetExpr(expr *SetExpression) interface{} {
    return nil
}

func (a *AstPrinter) visitThisExpr(expr *ThisExpr) interface{} {
    return nil 
}

func (a *AstPrinter) visitSuperExpr(expr *SuperExpr) interface{} {
    return nil 
}

// Helper methods
func (a *AstPrinter) printExpr(expr Expr) string {
    if expr == nil {
        return "nil"
    }
    result := expr.Accept(a)
    if result == nil {
        return "nil"
    }
    if s, ok := result.(string); ok {
        return s
    }
    return fmt.Sprintf("%v", result)
}

func (a *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(name)
	for _, expr := range exprs {
		sb.WriteString(" ")
		sb.WriteString(a.printExpr(expr))
	}
	sb.WriteString(")")
	return sb.String()
}
