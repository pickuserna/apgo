// Package apast defines a simplified AST format that is easy to interpret.
// We make a number of simplifying assumptions:
// * We assume that the code already compiles, so we don't check things like
//   type errors.
// * There are no operators; they are replaced by function calls.
package apast

import (
	"reflect"
)

type Stmt interface {
	apstmtNode()
}

type ExprStmt struct {
	E Expr
}

type AssignStmt struct {
	Lhs []Expr
	Rhs []Expr
}

type BlockStmt struct {
	Stmts []Stmt
}

func (*ExprStmt) apstmtNode() {}
func (*BlockStmt) apstmtNode() {}
func (*AssignStmt) apstmtNode() {}

type Expr interface {
	apexprNode()
}

type FuncCallExpr struct {
	Func Expr
	Args []Expr
}

type IdentExpr struct {
	Name string
}

type LiteralExpr struct {
	Val reflect.Value
}

func (*FuncCallExpr) apexprNode() {}
func (*IdentExpr) apexprNode() {}
func (*LiteralExpr) apexprNode() {}
