package main

type Stmt interface {
	Accept(visitor StmtVisitor) interface{}
}

type StmtVisitor interface {
	visitExpressionStmt(stmt *ExpressionStatement) interface{}
	visitPrintStmt(stmt *PrintStatement) interface{}
	visitVarStmt(stmt *VarStatement) interface{}
	visitBlockStmt(stmt *Block) interface{}
	visitIfStmt(stmt *IfStatement) interface{}
	visitWhileStmt(stmt *WhileStatement) interface{}
	visitFunctionStmt(stmt *FunctionStatement) interface{}
	visitReturnStmt(stmt *ReturnStatement) interface{}
	visitClassStmt(stmt *ClassStatement) interface{} 
}

type ExpressionStatement struct {
	Expression Expr
}

type PrintStatement struct {
	Value Expr
}

type VarStatement struct {
	Name        Token
	Initializer Expr
}

type Block struct {
	Statements []Stmt
}

type IfStatement struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

type WhileStatement struct {
	Condition Expr
	Body      Stmt
}

type FunctionStatement struct {
	Body   []Stmt
	Params []Token
	Name   Token
}

type ReturnStatement struct {
	Value   Expr
	Keyword Token
}

type ClassStatement struct {
	Name Token 
	Methods []*FunctionStatement
	Superclass *VariableExpr
}

func (s *ExpressionStatement) Accept(visitor StmtVisitor) interface{} {
    return visitor.visitExpressionStmt(s)
}

func (s *PrintStatement) Accept(visitor StmtVisitor) interface{} {
	return visitor.visitPrintStmt(s)
}

func (s *VarStatement) Accept(visitor StmtVisitor) interface{} {
	return visitor.visitVarStmt(s)
}

func (b *Block) Accept(visitor StmtVisitor) interface{} {
	visitor.visitBlockStmt(b)
	return nil
}

func (s *IfStatement) Accept(visitor StmtVisitor) interface{} {
	return visitor.visitIfStmt(s)
}

func (s *WhileStatement) Accept(visitor StmtVisitor) interface{} {
	return visitor.visitWhileStmt(s)
}

func (s *FunctionStatement) Accept(visitor StmtVisitor) interface{} {
	return visitor.visitFunctionStmt(s)
}

func (s *ReturnStatement) Accept(visitor StmtVisitor) interface{} {
	return visitor.visitReturnStmt(s)
}

func (s *ClassStatement) Accept(visitor StmtVisitor) interface{} {
	return visitor.visitClassStmt(s)
}
