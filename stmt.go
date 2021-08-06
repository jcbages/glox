package main

type Stmt interface {
	accept(visitor StmtVisitor) error
}

type StmtBlock struct {
	Statements []Stmt
}

type StmtVarDeclaration struct {
	Name        Token
	Initializer Expr
}

type StmtExpression struct {
	Expression Expr
}

type StmtPrint struct {
	Expression Expr
}

func (stmt StmtBlock) accept(visitor StmtVisitor) error {
	return visitor.VisitStmtBlock(stmt)
}

func (stmt StmtVarDeclaration) accept(visitor StmtVisitor) error {
	return visitor.VisitStmtVarDeclaration(stmt)
}

func (stmt StmtExpression) accept(visitor StmtVisitor) error {
	return visitor.VisitStmtExpression(stmt)
}

func (stmt StmtPrint) accept(visitor StmtVisitor) error {
	return visitor.VisitStmtPrint(stmt)
}
