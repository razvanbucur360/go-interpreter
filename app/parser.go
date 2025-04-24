package main

import (
	"fmt"
)

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	p := new(Parser)
	p.tokens = tokens
	p.current = 0
	return p
}

func (p *Parser) statement() Stmt {
	if p.match(PRINT) {
		return p.printStatement()
	}

	if p.match(FOR) {
		return p.forStatement()
	}

	if p.match(WHILE) {
		return p.whileStatement()
	}

	if p.match(RETURN) {
		return p.returnStatement()
	}

	if p.match(LEFT_BRACE) {
		return &Block{
			Statements: p.block(),
		}
	}

	if p.match(IF) {
		return p.ifStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) block() []Stmt {
	statements := make([]Stmt, 0)
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	p.consume(RIGHT_BRACE, "Expect '}' after block.")
	return statements
}

func (p *Parser) forStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'for'")

	var initializer Stmt
	if p.match(SEMICOLON) {
		initializer = nil
	} else if p.match(VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	if hasError {
		return &ExpressionStatement{Expression: nil}
	}

	var condition Expr
	if !p.check(SEMICOLON) {
		condition = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after loop condition")

	if hasError {
		return &ExpressionStatement{Expression: nil}
	}

	var increment Expr
	if !p.check(RIGHT_PAREN) {
		increment = p.expression()
	}
	p.consume(RIGHT_PAREN, "Expect ')' after for clauses")

	body := p.statement()

	if increment != nil {
		body = &Block{
			Statements: []Stmt{
				body,
				&ExpressionStatement{
					Expression: increment,
				},
			},
		}
	}

	if condition == nil {
		condition = &LiteralExpr{
			Value: true,
		}
	}

	body = &WhileStatement{
		Condition: condition,
		Body:      body,
	}

	if initializer != nil {
		body = &Block{
			Statements: []Stmt{
				initializer,
				body,
			},
		}
	}

	return body
}

func (p *Parser) whileStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'while'")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after 'while'")

	body := p.statement()

	return &WhileStatement{
		Condition: condition,
		Body:      body,
	}
}

func (p *Parser) returnStatement() Stmt {
	keyword := p.previous()
	var value Expr = nil
	if !p.check(SEMICOLON) {
		value = p.expression()
	}

	p.consume(SEMICOLON, "Expect ';' after return value.")

	return &ReturnStatement{
		Value:   value,
		Keyword: keyword,
	}
}

func (p *Parser) ifStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'if'")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after if condition")

	thenBranch := p.statement()
	var elseBranch Stmt = nil
	if p.match(ELSE) {
		elseBranch = p.statement()
	}

	return &IfStatement{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()

	if !p.check(SEMICOLON) && p.peek().TokenType != EOF {
		handleRuntimeError(RuntimeError{
			Token:   p.peek(),
			Message: "Expect ';' after value",
		})
	} else if p.check(SEMICOLON) {
		p.advance()
	}

	return &PrintStatement{
		Value: value,
	}
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()

	if !p.check(SEMICOLON) && p.peek().TokenType != EOF {
		handleRuntimeError(RuntimeError{
			Token:   p.peek(),
			Message: "Expect ';' after value",
		})
	} else if p.check(SEMICOLON) {
		p.advance()
	}

	return &ExpressionStatement{
		Expression: expr,
	}
}

func (p *Parser) varDeclaration() Stmt {
	name := p.consume(IDENTIFIER, "Expect variable name.")

	var initializer Expr
	if p.match(EQUAL) {
		initializer = p.expression()
	}

	p.consume(SEMICOLON, "Expect ';' after variable declaration.")

	return &VarStatement{
		Name:        name,
		Initializer: initializer,
	}
}

func (p *Parser) funDeclaration(kind string) *FunctionStatement {
	name := p.consume(IDENTIFIER, fmt.Sprintf("Expect %s name.", kind))

	p.consume(LEFT_PAREN, fmt.Sprintf("Expect '(' after %s name.", kind))

	parameters := make([]Token, 0)

	if !p.check(RIGHT_PAREN) {
		for {
			if len(parameters) > 255 {
				p.error(p.peek(), "Cannot have more than 255 parameters.")
			}
			parameters = append(parameters, p.consume(IDENTIFIER, "Expect parameter name."))
			if !p.match(COMMA) {
				break
			}
		}
	}

	p.consume(RIGHT_PAREN, "Expect ')' after parameters.")

	p.consume(LEFT_BRACE, fmt.Sprintf("Expect '{' before %s body", kind))

	body := p.block()

	return &FunctionStatement{
		Body:   body,
		Name:   name,
		Params: parameters,
	}
}

func (p *Parser) classDeclaration(kind string) Stmt {
	name := p.consume(IDENTIFIER, fmt.Sprintf("Expect %s name.", kind))

	var superclass *VariableExpr

	if p.match(LESS) {
		p.consume(IDENTIFIER, "Expect superclass name.")
		superclass = &VariableExpr{
			Name: p.previous(),
		}
	}

	p.consume(LEFT_BRACE, fmt.Sprintf("Expect '{' after %s name.", kind))

	methods := make([]*FunctionStatement, 0)

	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		methods = append(methods, p.funDeclaration("method"))
	}

	p.consume(RIGHT_BRACE, "Expect '}' after class body.")

	return &ClassStatement{
		Name: name,
		Methods: methods,
		Superclass: superclass,
	}
}

func (p *Parser) declaration() Stmt {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(*ParseError); ok {
				p.synchronize()
			} else {
				panic(r) // Re-panic for unexpected errors
			}
		}
	}()

	if p.match(VAR) {
		return p.varDeclaration()
	}

	if p.match(FUN) {
		return p.funDeclaration("function")
	}

	if p.match(CLASS) {
		return p.classDeclaration("class")
	}

	return p.statement()
}

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) assignment() Expr {
	expr := p.or()

	if p.match(EQUAL) {
		equal := p.previous()
		value := p.assignment()

		if varExpr, ok := expr.(*VariableExpr); ok {
			name := varExpr.Name
			return &AssignmentExpr{
				Name:  name,
				Value: value,
			}
		} else if getExpr, ok := expr.(*GetExpression); ok {
			name := getExpr.Name
			object := getExpr.Object
			return &SetExpression{
				Value: value,
				Name: name, 
				Object: object,
			}
		}

		p.error(equal, "Invalid assignment target.")
	}

	return expr
}

func (p *Parser) or() Expr {
	expr := p.and()
	for p.match(OR) {
		operator := p.previous()
		right := p.and()
		expr = &LogicalExpr{
			Operator: operator,
			Left:     expr,
			Right:    right,
		}
	}
	return expr
}

func (p *Parser) and() Expr {
	expr := p.equality()
	for p.match(AND) {
		operator := p.previous()
		right := p.equality()
		expr = &LogicalExpr{
			Operator: operator,
			Left:     expr,
			Right:    right,
		}
	}
	return expr
}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = &BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(LESS, LESS_EQUAL, GREATER, GREATER_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = &BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = &BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()
		expr = &BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}

func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()
		expr := &UnaryExpr{
			Operator: operator,
			Right:    right,
		}
		return expr
	}

	return p.call()
}

func (p *Parser) call() Expr {
	expr := p.primary()
	for true {
		if p.match(LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else if p.match(DOT) {
			name := p.consume(IDENTIFIER, "Expect property name after '.'.")
			expr = &GetExpression{
				Name: name,
				Object: expr,
			}
		} else {
			break 
		}
	}
	return expr
}

func (p *Parser) primary() Expr {
	if p.match(FALSE) {
		return &LiteralExpr{
			Value: false,
		}
	}

	if p.match(THIS) {
		return &ThisExpr{
			Keyword: p.previous(),
		}
	}

	if p.match(TRUE) {
		return &LiteralExpr{
			Value: true,
		}
	}

	if p.match(NIL) {
		return &LiteralExpr{
			Value: nil,
		}
	}

	if p.match(NUMBER, STRING) {
		return &LiteralExpr{
			p.previous().Literal,
		}
	}

	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return &GroupingExpr{
			Expression: expr,
		}
	}

	if p.match(IDENTIFIER) {
		return &VariableExpr{
			Name: p.previous(),
		}
	}

	if p.match(SUPER) {
		keyword := p.previous()
		p.consume(DOT, "Expect '.' after 'super'.")
		method := p.consume(IDENTIFIER, "Expect superclass method name.")
		return &SuperExpr{
			Keyword: keyword,
			Method: method,
		}
	}

	panic(p.error(p.peek(), "Expect expression."))
}

func (p *Parser) match(tokens ...TokenType) bool {
	for _, token := range tokens {
		if p.check(token) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(token TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.tokens[p.current].TokenType == token
}

func (p *Parser) isAtEnd() bool {
	if p.peek().TokenType == EOF {
		return true
	}
	return false
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) consume(TokenType TokenType, message string) Token {
	if p.check(TokenType) {
		return p.advance()
	}

	panic(p.error(p.peek(), message))
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().TokenType == SEMICOLON {
			return
		}

		switch p.peek().TokenType {
		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT, RETURN:
			return
		}

		p.advance()
	}
}

func (p *Parser) finishCall(callee Expr) Expr {
	arguments := make([]Expr, 0)

	if !p.check(RIGHT_PAREN) {
		for {
			if len(arguments) > 255 {
				p.error(p.peek(), "Can't have more than 255 arguments.")
			}
			arguments = append(arguments, p.expression())
			if !p.match(COMMA) {
				break
			}
		}
	}

	paren := p.consume(RIGHT_PAREN, "Expect ')' after arguments.")

	return &CallExpression{
		Callee:      callee,
		Arguments:   arguments,
		Parenthesis: paren,
	}
}

func (p *Parser) parse() []Stmt {
	statements := make([]Stmt, 0)

	for !p.isAtEnd() {
		defer func() {
			if r := recover(); r != nil {
				if _, ok := r.(*ParseError); ok {
					p.synchronize()
				} else {
					panic(r)
				}
			}
		}()

		declaration := p.declaration()
		if declaration != nil {
			statements = append(statements, declaration)
		}
	}

	return statements
}

type ParseError struct {
	Token   Token
	Message string
}

func (e *ParseError) Error() string {
	return e.Message
}

func (p *Parser) error(token Token, message string) *ParseError {
	if token.TokenType == EOF {
		report(token.Line, " at end", message)
	} else {
		report(token.Line, fmt.Sprintf(" at '%s'", token.Lexeme), message)
	}

	return &ParseError{
		Token:   token,
		Message: message,
	}
}
