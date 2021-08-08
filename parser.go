package main

import (
	"errors"
	"fmt"
)

type Parser struct {
	Tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		Tokens:  tokens,
		current: 0,
	}
}

func (parser *Parser) Parse() ([]Stmt, error) {
	var statements []Stmt
	for !parser.isAtEnd() {
		fmt.Printf("NEXT TOKEN %v\n", parser.peek())
		stmt, err := parser.declaration()
		if err != nil {
			return nil, err
		}

		statements = append(statements, stmt)
	}

	return statements, nil
}

func (parser *Parser) declaration() (Stmt, error) {
	if parser.match(VAR) {
		stmt, err := parser.varDeclarationStatement()
		if err != nil {
			parser.synchronize()
			return nil, err
		} else {
			return stmt, nil
		}
	}

	stmt, err := parser.statement()
	if err != nil {
		parser.synchronize()
		return nil, err
	} else {
		return stmt, nil
	}
}

func (parser *Parser) varDeclarationStatement() (Stmt, error) {
	name, err := parser.consume(IDENTIFIER, "Expected variable name")
	if err != nil {
		return nil, err
	}

	var initializer Expr
	if parser.match(EQUAL) {
		initializer, err = parser.expression()
		if err != nil {
			return nil, err
		}
	}

	if _, err := parser.consume(SEMICOLON, "Expected ';' after value"); err != nil {
		return nil, err
	}

	return StmtVarDeclaration{
		Name:        name,
		Initializer: initializer,
	}, nil
}

func (parser *Parser) block() ([]Stmt, error) {
	var statements []Stmt

	for !parser.check(RIGHT_BRACE) && !parser.isAtEnd() {
		stmt, err := parser.declaration()
		if err != nil {
			return nil, err
		} else {
			statements = append(statements, stmt)
		}
	}

	if _, err := parser.consume(RIGHT_BRACE, "Expect '}' after block"); err != nil {
		return nil, err
	}

	return statements, nil
}

func (parser *Parser) statement() (Stmt, error) {
	if parser.match(FOR) {
		return parser.forStatement()
	} else if parser.match(WHILE) {
		return parser.whileStatement()
	} else if parser.match(IF) {
		return parser.ifStatement()
	} else if parser.match(PRINT) {
		return parser.printStatement()
	} else if parser.match(LEFT_BRACE) {
		statements, err := parser.block()
		if err != nil {
			return nil, err
		} else {
			return StmtBlock{Statements: statements}, nil
		}
	} else {
		return parser.expressionStatement()
	}
}

func (parser *Parser) forStatement() (Stmt, error) {
	var err error
	if _, err := parser.consume(LEFT_PAREN, "Expected '(' after while"); err != nil {
		return nil, err
	}

	var initializer Stmt
	if parser.match(VAR) {
		initializer, err = parser.varDeclarationStatement()
		if err != nil {
			return nil, err
		}
	} else if !parser.match(SEMICOLON) {
		initializer, err = parser.expressionStatement()
		if err != nil {
			return nil, err
		}
	}

	var condition Expr
	if !parser.match(SEMICOLON) {
		condition, err = parser.expression()
		if err != nil {
			return nil, err
		}
	}

	if _, err := parser.consume(SEMICOLON, "Expected ';' after condition"); err != nil {
		return nil, err
	}

	var increment Expr
	if !parser.match(RIGHT_PAREN) {
		increment, err = parser.expression()
		if err != nil {
			return nil, err
		}
	}

	if _, err := parser.consume(SEMICOLON, "Expected ')' after increment"); err != nil {
		return nil, err
	}

	body, err := parser.statement()
	if err != nil {
		return nil, err
	}

	// De-sugarize the for-loop into a while-loop

	if increment != nil {
		body = StmtBlock{
			Statements: []Stmt{
				body,
				StmtExpression{Expression: increment},
			},
		}
	}

	if condition == nil {
		condition = ExprLiteral{Value: true}
	}
	body = StmtWhile{
		Condition: condition,
		Body:      body,
	}

	if initializer != nil {
		body = StmtBlock{
			Statements: []Stmt{
				initializer,
				body,
			},
		}
	}

	return body, nil
}

func (parser *Parser) whileStatement() (Stmt, error) {
	if _, err := parser.consume(LEFT_PAREN, "Expected '(' after while"); err != nil {
		return nil, err
	}

	condition, err := parser.expression()
	if err != nil {
		return nil, err
	}

	if _, err := parser.consume(RIGHT_PAREN, "Expected ')' after expression"); err != nil {
		return nil, err
	}

	body, err := parser.statement()
	if err != nil {
		return nil, err
	}

	return StmtWhile{
		Condition: condition,
		Body:      body,
	}, nil
}

func (parser *Parser) ifStatement() (Stmt, error) {
	if _, err := parser.consume(LEFT_PAREN, "Expected '(' after if"); err != nil {
		return nil, err
	}

	condition, err := parser.expression()
	if err != nil {
		return nil, err
	}

	if _, err := parser.consume(RIGHT_PAREN, "Expected ')' after expression"); err != nil {
		return nil, err
	}

	thenBranch, err := parser.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch Stmt
	if parser.match(ELSE) {
		elseBranch, err = parser.statement()
		if err != nil {
			return nil, err
		}
	}

	return StmtIf{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}, nil
}

func (parser *Parser) printStatement() (Stmt, error) {
	expr, err := parser.expression()
	if err != nil {
		return nil, err
	}

	if _, err := parser.consume(SEMICOLON, "Expected ';' after value"); err != nil {
		return nil, err
	}

	return StmtPrint{Expression: expr}, nil
}

func (parser *Parser) expressionStatement() (Stmt, error) {
	expr, err := parser.expression()
	if err != nil {
		return nil, err
	}

	if _, err := parser.consume(SEMICOLON, "Expected ';' after expression"); err != nil {
		return nil, err
	}

	return StmtExpression{Expression: expr}, nil
}

func (parser *Parser) expression() (Expr, error) {
	return parser.comma()
}

func (parser *Parser) binary(operand func() (Expr, error), isLogical bool, tokenTypes ...TokenType) (Expr, error) {
	expr, err := operand()
	if err != nil {
		return nil, err
	}

	for parser.match(tokenTypes...) {
		operator := parser.previous()
		right, err := operand()
		if err != nil {
			return nil, err
		}

		if isLogical {
			expr = ExprLogical{
				Operator: operator,
				Left:     expr,
				Right:    right,
			}
		} else {
			expr = ExprBinary{
				Operator: operator,
				Left:     expr,
				Right:    right,
			}
		}
	}

	return expr, nil
}

func (parser *Parser) match(tokenTypes ...TokenType) bool {
	for _, tokenType := range tokenTypes {
		if parser.check(tokenType) {
			parser.advance()
			return true
		}
	}

	return false
}

func (parser *Parser) check(tokenType TokenType) bool {
	if parser.isAtEnd() {
		return false
	} else {
		return parser.Tokens[parser.current].TokenType == tokenType
	}
}

func (parser *Parser) isAtEnd() bool {
	return parser.Tokens[parser.current].TokenType == EOF
}

func (parser *Parser) advance() Token {
	parser.current += 1
	return parser.previous()
}

func (parser *Parser) peek() Token {
	return parser.Tokens[parser.current]
}

func (parser *Parser) previous() Token {
	return parser.Tokens[parser.current-1]
}

func (parser *Parser) comma() (Expr, error) {
	return parser.binary(parser.assignment, false, COMMA)
}

func (parser *Parser) assignment() (Expr, error) {
	expr, err := parser.or()
	if err != nil {
		return nil, err
	}

	if parser.match(EQUAL) {
		equals := parser.previous()
		value, err := parser.assignment()
		if err != nil {
			return nil, err
		}

		if varExpr, ok := expr.(ExprVariable); ok {
			return ExprAssign{
				Name:  varExpr.Name,
				Value: value,
			}, nil
		} else {
			LoxTokenError(equals, "Invalid assignment target")
			return nil, errors.New("Invalid assignment target")
		}
	} else {
		return expr, nil
	}
}

func (parser *Parser) or() (Expr, error) {
	return parser.binary(parser.and, true, OR)
}

func (parser *Parser) and() (Expr, error) {
	return parser.binary(parser.equality, true, AND)
}

func (parser *Parser) equality() (Expr, error) {
	return parser.binary(parser.comparison, false, BANG_EQUAL, EQUAL_EQUAL)
}

func (parser *Parser) comparison() (Expr, error) {
	return parser.binary(parser.term, false, GREATER, GREATER_EQUAL, LESS, LESS_EQUAL)
}

func (parser *Parser) term() (Expr, error) {
	return parser.binary(parser.factor, false, MINUS, PLUS)
}

func (parser *Parser) factor() (Expr, error) {
	return parser.binary(parser.unary, false, STAR, SLASH)
}

func (parser *Parser) unary() (Expr, error) {
	if parser.match(BANG, MINUS) {
		operator := parser.previous()
		right, err := parser.unary()
		return ExprUnary{
			Operator: operator,
			Right:    right,
		}, err
	} else {
		return parser.primary()
	}
}

func (parser *Parser) primary() (Expr, error) {
	switch true {
	case parser.match(FALSE):
		return ExprLiteral{Value: false}, nil
	case parser.match(TRUE):
		return ExprLiteral{Value: true}, nil
	case parser.match(NIL):
		return ExprLiteral{Value: nil}, nil
	case parser.match(NUMBER, STRING):
		return ExprLiteral{Value: parser.previous().Literal}, nil
	case parser.match(IDENTIFIER):
		return ExprVariable{Name: parser.previous()}, nil
	case parser.match(LEFT_PAREN):
		expr, err := parser.expression()
		if err != nil {
			return nil, err
		}

		if _, err = parser.consume(RIGHT_PAREN, "Expected ')' after expression"); err != nil {
			return nil, err
		}

		return ExprGrouping{Expression: expr}, err
	default:
		LoxTokenError(parser.peek(), "Expected expression")
		return nil, errors.New("Expected expression")
	}
}

func (parser *Parser) consume(tokenType TokenType, message string) (Token, error) {
	if parser.check(tokenType) {
		return parser.advance(), nil
	} else {
		LoxTokenError(parser.peek(), message)
		return Token{}, errors.New("Unexpected token")
	}
}

func (parser *Parser) synchronize() {
	parser.advance()
	for !parser.isAtEnd() {
		if parser.previous().TokenType == SEMICOLON {
			return
		}

		switch parser.peek().TokenType {
		case CLASS, FOR, FUN, IF, PRINT, RETURN, VAR, WHILE:
			return
		}

		parser.advance()
	}
}
