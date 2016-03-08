package apevaluator

import (
	"github.com/alangpierce/apgo/apruntime"
	"github.com/alangpierce/apgo/apast"
	"reflect"
	"fmt"
)

func EvaluateStmt(ctx apruntime.Context, stmt apast.Stmt) {
	switch stmt := stmt.(type) {
	case *apast.ExprStmt:
		evaluateExpr(ctx, stmt.E)
	case *apast.BlockStmt:
		for _, line := range stmt.Stmts {
			EvaluateStmt(ctx, line)
		}
	case *apast.AssignStmt:
		if len(stmt.Lhs) != len(stmt.Rhs) {
			panic("Multiple assign with differing lengths not implemented.")
		}
		values := []reflect.Value{}
		for _, rhsExpr := range stmt.Rhs {
			values = append(values, evaluateExpr(ctx, rhsExpr))
		}
		for i, value := range values {
			lvalue := stmt.Lhs[i]
			if lvalue, ok := lvalue.(*apast.IdentExpr); ok {
				ctx[lvalue.Name] = value
			}
		}
	case *apast.EmptyStmt:
		// Do nothing.
	case *apast.IfStmt:
		// TODO: Handle scopes properly, if necessary.
		EvaluateStmt(ctx, stmt.Init)
		condValue := evaluateExpr(ctx, stmt.Cond)
		if condValue.Interface().(bool) {
			EvaluateStmt(ctx, stmt.Body)
		} else {
			EvaluateStmt(ctx, stmt.Else)
		}
	default:
		panic(fmt.Sprint("Statement eval not implemented: ", reflect.TypeOf(stmt)))
	}
}


func evaluateExpr(ctx apruntime.Context, expr apast.Expr) reflect.Value {
	switch expr := expr.(type) {
	case *apast.FuncCallExpr:
		funcValue := evaluateExpr(ctx, expr.Func)
		argValues := []reflect.Value{}
		for _, argExpr := range expr.Args {
			argValues = append(argValues, evaluateExpr(ctx, argExpr))
		}
		// TODO: Handle multiple return values.
		return funcValue.Call(argValues)[0]
	case *apast.IdentExpr:
		return ctx[expr.Name]
	case *apast.LiteralExpr:
		return expr.Val
	default:
		panic(fmt.Sprint("Expression eval not implemented: ", reflect.TypeOf(expr)))
	}
}