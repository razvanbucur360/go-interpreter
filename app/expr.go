package main

type Expr interface {
	Accept(visitor ExprVisitor) interface{}
}

type ExprVisitor interface {
	visitBinaryExpr(expr *BinaryExpr) interface{}
	visitUnaryExpr(expr *UnaryExpr) interface{}
	visitGroupingExpr(expr *GroupingExpr) interface{}
	visitLiteralExpr(expr *LiteralExpr) interface{}
	visitVariableExpr(expr *VariableExpr) interface{}
	visitAssignmentExpr(expr *AssignmentExpr) interface{}
	visitLogicalExpr(expr *LogicalExpr) interface{} 
	visitCallExpr(expr *CallExpression) interface{}
	visitGetExpr(expr *GetExpression) interface{} 
	visitSetExpr(expr *SetExpression) interface{} 
	visitThisExpr(expr *ThisExpr) interface{}
	visitSuperExpr(expr *SuperExpr) interface{} 
}

type BinaryExpr struct {
	Left    Expr
	Operator Token
	Right    Expr
}

type UnaryExpr struct {
	Operator Token
	Right    Expr
}

type GroupingExpr struct {
	Expression Expr
}

type LiteralExpr struct {
	Value interface{}
}

type LogicalExpr struct {
	Operator Token 
	Left Expr 
	Right Expr 
}

type VariableExpr struct {
	Name Token
}

type AssignmentExpr struct {
	Name Token 
	Value Expr 
}

type CallExpression struct {
	Callee Expr 
	Arguments []Expr 
	Parenthesis Token 
}

type GetExpression struct {
	Name Token 
	Object Expr
}

type SetExpression struct {
	Name Token 
	Object Expr 
	Value Expr 
}

type ThisExpr struct {
	Keyword Token 
}

type SuperExpr struct {
	Keyword Token 
	Method Token 
}

func (e *BinaryExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.visitBinaryExpr(e)
}

func (e *UnaryExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.visitUnaryExpr(e)
}

func (e *GroupingExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.visitGroupingExpr(e)
}

func (e *LiteralExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.visitLiteralExpr(e)
}

func (e *VariableExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.visitVariableExpr(e)
}

func (e *AssignmentExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.visitAssignmentExpr(e)
}

func (e *LogicalExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.visitLogicalExpr(e)
}

func (e *CallExpression) Accept(visitor ExprVisitor) interface{} {
	return visitor.visitCallExpr(e)
}

func (e *GetExpression) Accept(visitor ExprVisitor) interface{} {
	return visitor.visitGetExpr(e)
}

func (e *SetExpression) Accept(visitor ExprVisitor) interface{} {
	return visitor.visitSetExpr(e)
}

func (e *ThisExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.visitThisExpr(e)
}

func (e *SuperExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.visitSuperExpr(e)
}
/*
Visitor:
1. Parent class which has the accept(visitor GeneralVisitor)
2. Child classes which inherit from Parent class
3. GeneralVisitor interface which declares methods for each child class. e.g. visitBinary(expr BinaryExpr) visitUnary (expr UnaryExpr)
4. Child classes call the method corresponding to their class type. e.g. ExprBinary calls visitBinary(this)
5. To extend functionality, declare a class that implements the GenearalVisitor interface and override the method.
6. Declare a object as GeneralVisitor but instantiate it as the class which implements the specialized behavior.
*/