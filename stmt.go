package main

type Stmt interface {
	accept(visitor StmtVisitor) error
}

type StmtFunction struct {
	Name       Token
	Parameters []Token
	Body       []Stmt
}

type StmtWhile struct {
	Condition Expr
	Body      Stmt
}

type StmtIf struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
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

type StmtReturn struct {
	Keyword    Token
	Expression Expr
}

func (stmt StmtFunction) accept(visitor StmtVisitor) error {
	return visitor.VisitStmtFunction(stmt)
}

func (stmt StmtWhile) accept(visitor StmtVisitor) error {
	return visitor.VisitStmtWhile(stmt)
}

func (stmt StmtIf) accept(visitor StmtVisitor) error {
	return visitor.VisitStmtIf(stmt)
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

func (stmt StmtReturn) accept(visitor StmtVisitor) error {
	return visitor.VisitStmtReturn(stmt)
}
