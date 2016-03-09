// Package apast defines a simplified AST format that is easy to interpret.
// We make a number of simplifying assumptions:
// * We assume that the code already compiles, so we don't check things like
//   type errors.
// * There are no operators; they are replaced by function calls.
package apast

import (
	"reflect"
	"fmt"
)

type Package struct {
	Funcs map[string]*FuncDecl
}

type FuncDecl struct {
	Body       Stmt
	ParamNames []string
}

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

type EmptyStmt struct {
}

// All fields are required.
type IfStmt struct {
	Init Stmt
	Cond Expr
	Body Stmt
	Else Stmt
}

type ReturnStmt struct {
	Results []Expr
}

func (*ExprStmt) apstmtNode() {}
func (*AssignStmt) apstmtNode() {}
func (*BlockStmt) apstmtNode() {}
func (*EmptyStmt) apstmtNode() {}
func (*IfStmt) apstmtNode() {}
func (*ReturnStmt) apstmtNode() {}

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


func (e *FuncCallExpr) String() string {
	return fmt.Sprintf("FuncCall{%s,%s}", e.Func, e.Args)
}
func (e *IdentExpr) String() string {
	return fmt.Sprintf("Ident{%s}", e.Name)
}
func (e *LiteralExpr) String() string {
	return fmt.Sprintf("Literal{%s}", e.Val.Interface())
}