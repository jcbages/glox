package main

type ExprVisitor interface {
	VisitExprBinary(expr ExprBinary) (interface{}, error)
	VisitExprGrouping(expr ExprGrouping) (interface{}, error)
	VisitExprLiteral(expr ExprLiteral) (interface{}, error)
	VisitExprVariable(expr ExprVariable) (interface{}, error)
	VisitExprUnary(expr ExprUnary) (interface{}, error)
	VisitExprAssign(expr ExprAssign) (interface{}, error)
}

type StmtVisitor interface {
	VisitStmtVarDeclaration(stmt StmtVarDeclaration) error
	VisitStmtExpression(stmt StmtExpression) error
	VisitStmtPrint(stmt StmtPrint) error
	VisitStmtBlock(stmt StmtBlock) error
}
