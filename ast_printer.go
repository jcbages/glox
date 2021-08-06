package main

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

func (ast *AstPrinter) print(statements []Stmt) string {
	var sb strings.Builder
	sb.WriteString("(program")

	for _, stmt := range statements {
		if stmtExpr, ok := stmt.(StmtExpression); ok {
			value, _ := stmtExpr.Expression.accept(ast)
			sb.WriteString(fmt.Sprintf(" (%v)", value))
		} else if stmtPrint, ok := stmt.(StmtPrint); ok {
			value, _ := stmtPrint.Expression.accept(ast)
			sb.WriteString(fmt.Sprintf(" (print (%v))", value))
		} else {
			sb.WriteString("statement without expression")
		}
	}

	sb.WriteString(")")
	return sb.String()
}

func (ast *AstPrinter) VisitExprBinary(expr ExprBinary) (interface{}, error) {
	return ast.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right), nil
}

func (ast *AstPrinter) VisitExprGrouping(expr ExprGrouping) (interface{}, error) {
	return ast.parenthesize("group", expr.Expression), nil
}

func (ast *AstPrinter) VisitExprLiteral(expr ExprLiteral) (interface{}, error) {
	if expr.Value == nil {
		return "nil", nil
	} else {
		return expr.Value, nil
	}
}

func (ast *AstPrinter) VisitExprAssign(expr ExprAssign) (interface{}, error) {
	return fmt.Sprintf("(= %v %v)", expr.Name.Lexeme, expr.Value), nil
}

func (ast *AstPrinter) VisitExprVariable(expr ExprVariable) (interface{}, error) {
	return expr.Name.Lexeme, nil
}

func (ast *AstPrinter) VisitExprUnary(expr ExprUnary) (interface{}, error) {
	return ast.parenthesize(expr.Operator.Lexeme, expr.Right), nil
}

func (ast *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	var sb strings.Builder

	sb.WriteString("(")
	sb.WriteString(name)
	for _, expr := range exprs {
		sb.WriteString(" ")
		value, _ := expr.accept(ast)
		sb.WriteString(fmt.Sprintf("%v", value))
	}
	sb.WriteString(")")

	return sb.String()
}
