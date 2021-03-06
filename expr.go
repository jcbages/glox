package main

type Expr interface {
	accept(visitor ExprVisitor) (interface{}, error)
}

type ExprCall struct {
	Callee    Expr
	Paren     Token
	Arguments []Expr
}

type ExprBinary struct {
	Operator Token
	Left     Expr
	Right    Expr
}

type ExprLogical struct {
	Operator Token
	Left     Expr
	Right    Expr
}

type ExprGrouping struct {
	Expression Expr
}

type ExprAssign struct {
	Name  Token
	Value Expr
}

type ExprLiteral struct {
	Value interface{}
}

type ExprVariable struct {
	Name Token
}

type ExprUnary struct {
	Operator Token
	Right    Expr
}

func (expr ExprCall) accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitExprCall(expr)
}

func (expr ExprLogical) accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitExprLogical(expr)
}

func (expr ExprAssign) accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitExprAssign(expr)
}

func (expr ExprBinary) accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitExprBinary(expr)
}

func (expr ExprGrouping) accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitExprGrouping(expr)
}

func (expr ExprLiteral) accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitExprLiteral(expr)
}

func (expr ExprVariable) accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitExprVariable(expr)
}

func (expr ExprUnary) accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitExprUnary(expr)
}
