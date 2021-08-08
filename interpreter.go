package main

import (
	"fmt"
	"math"
)

const EPS = 1e-9

type RuntimeError struct {
	Token   Token
	Message string
}

func (err RuntimeError) Error() string {
	return err.Message
}

type Interpreter struct {
	environment *Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		environment: NewEnvironment(nil),
	}
}

func (intr *Interpreter) Interpret(statements []Stmt) error {
	for _, stmt := range statements {
		if err := intr.execute(stmt); err != nil {
			LoxRuntimeError(err.(RuntimeError).Token, err.Error())
		}
	}

	return nil
}

func (intr *Interpreter) stringify(value interface{}) string {
	if value == nil {
		return "nil"
	}

	if f, ok := value.(float64); ok && math.Abs(f-math.Round(f)) <= EPS {
		return fmt.Sprintf("%v", math.Round(f))
	}

	return fmt.Sprintf("%v", value)
}

func (intr *Interpreter) executeBlock(statements []Stmt, environment *Environment) error {
	previous := intr.environment
	intr.environment = environment

	for _, stmt := range statements {
		if err := intr.execute(stmt); err != nil {
			intr.environment = previous
			return err
		}
	}

	intr.environment = previous
	return nil
}

func (intr *Interpreter) execute(stmt Stmt) error {
	return stmt.accept(intr)
}

func (intr *Interpreter) evaluate(expr Expr) (interface{}, error) {
	return expr.accept(intr)
}

func (intr *Interpreter) VisitStmtWhile(stmt StmtWhile) error {
	for {
		value, err := intr.evaluate(stmt.Condition)
		if err != nil {
			return err
		}

		if !intr.isTruthy(value) {
			return nil
		}

		if err := intr.execute(stmt.Body); err != nil {
			return err
		}
	}
}

func (intr *Interpreter) VisitStmtIf(stmt StmtIf) error {
	value, err := intr.evaluate(stmt.Condition)

	if err != nil {
		return err
	} else if intr.isTruthy(value) {
		return intr.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		return intr.execute(stmt.ElseBranch)
	} else {
		return nil
	}
}

func (intr *Interpreter) VisitStmtBlock(stmt StmtBlock) error {
	return intr.executeBlock(stmt.Statements, NewEnvironment(intr.environment))
}

func (intr *Interpreter) VisitStmtVarDeclaration(stmt StmtVarDeclaration) error {
	if stmt.Initializer != nil {
		value, err := intr.evaluate(stmt.Initializer)
		if err != nil {
			return err
		} else {
			intr.environment.Define(stmt.Name.Lexeme, value)
		}
	} else {
		intr.environment.Define(stmt.Name.Lexeme, nil)
	}

	return nil
}

func (intr *Interpreter) VisitStmtPrint(stmt StmtPrint) error {
	value, err := intr.evaluate(stmt.Expression)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", value)
	return nil
}

func (intr *Interpreter) VisitStmtExpression(stmt StmtExpression) error {
	_, err := intr.evaluate(stmt.Expression)
	return err
}

func (intr *Interpreter) VisitExprLogical(expr ExprLogical) (interface{}, error) {
	left, err := intr.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	if expr.Operator.TokenType == OR && intr.isTruthy(left) {
		return left, nil
	}

	if expr.Operator.TokenType == AND && !intr.isTruthy(left) {
		return left, nil
	}

	return intr.evaluate(expr.Right)
}

func (intr *Interpreter) VisitExprAssign(expr ExprAssign) (interface{}, error) {
	value, err := intr.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}

	if err := intr.environment.Assign(expr.Name, value); err != nil {
		return nil, err
	} else {
		return value, nil
	}
}

func (intr *Interpreter) VisitExprBinary(expr ExprBinary) (interface{}, error) {
	left, err := intr.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	right, err := intr.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.TokenType {
	// Operations
	case STAR:
		if err := intr.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil
	case SLASH:
		if err := intr.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case MINUS:
		if err := intr.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case PLUS:
		f1, ok1 := left.(float64)
		f2, ok2 := right.(float64)
		if ok1 && ok2 {
			return f1 + f2, nil
		}

		s1, ok1 := left.(string)
		s2, ok2 := right.(string)
		if ok1 && ok2 {
			return s1 + s2, nil
		}
	// Comparisons
	case GREATER:
		if err := intr.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case GREATER_EQUAL:
		if err := intr.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case LESS:
		if err := intr.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case LESS_EQUAL:
		if err := intr.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	// Equality / Inequality
	case BANG_EQUAL:
		return !intr.isEqual(left, right), nil
	case EQUAL_EQUAL:
		return intr.isEqual(left, right), nil
	}

	return nil, RuntimeError{
		Token:   expr.Operator,
		Message: "Operands must be two numbers or two strings",
	}
}

func (intr *Interpreter) VisitExprGrouping(expr ExprGrouping) (interface{}, error) {
	return intr.evaluate(expr.Expression)
}

func (intr *Interpreter) VisitExprLiteral(expr ExprLiteral) (interface{}, error) {
	return expr.Value, nil
}

func (intr *Interpreter) VisitExprVariable(expr ExprVariable) (interface{}, error) {
	return intr.environment.Get(expr.Name)
}

func (intr *Interpreter) VisitExprUnary(expr ExprUnary) (interface{}, error) {
	right, err := intr.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.TokenType {
	case MINUS:
		if err := intr.checkNumberOperands(expr.Operator, right); err != nil {
			return nil, err
		}
		return -right.(float64), nil
	case BANG:
		return !intr.isTruthy(right), nil
	}

	return nil, RuntimeError{
		Token:   expr.Operator,
		Message: "Unexpected error interpreting unary expression",
	}
}

func (intr *Interpreter) isTruthy(value interface{}) bool {
	if value == nil {
		return false
	}

	if boolean, ok := value.(bool); ok {
		return boolean
	}

	return true
}

func (intr *Interpreter) isEqual(a interface{}, b interface{}) bool {
	return a == b
}

func (intr *Interpreter) checkNumberOperands(token Token, operands ...interface{}) error {
	for _, operand := range operands {
		if _, ok := operand.(float64); !ok {
			return RuntimeError{
				Token:   token,
				Message: "Operands must be two numbers",
			}
		}
	}

	return nil
}
