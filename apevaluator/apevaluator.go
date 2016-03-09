package apevaluator

import (
	"github.com/alangpierce/apgo/apast"
	"reflect"
	"fmt"
)

func EvaluateFunc(pack *apast.Package, funcAst *apast.FuncDecl, args ...interface{}) []reflect.Value {
	ctx := NewContext(pack)
	for i, argName := range funcAst.ParamNames {
		ctx.assignValue(argName, reflect.ValueOf(args[i]))
	}
	EvaluateStmt(ctx, funcAst.Body)
	return ctx.returnValues
}

func EvaluateStmt(ctx *Context, stmt apast.Stmt) {
	switch stmt := stmt.(type) {
	case *apast.ExprStmt:
		evaluateExpr(ctx, stmt.E)
	case *apast.BlockStmt:
		for _, line := range stmt.Stmts {
			EvaluateStmt(ctx, line)
			// If this sub-statement returned, we don't want to
			// continue any further.
			if ctx.returnValues != nil || ctx.shouldBreak {
				return
			}
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
				ctx.assignValue(lvalue.Name, value)
			} else {
				panic("Only assignment to identifiers supported for now.")
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
	case *apast.ForStmt:
		// TODO: Handle scopes properly, if necessary.
		EvaluateStmt(ctx, stmt.Init)
		for {
			condValue := evaluateExpr(ctx, stmt.Cond)
			if !condValue.Interface().(bool) {
				break
			}
			EvaluateStmt(ctx, stmt.Body)
			if ctx.shouldBreak {
				ctx.shouldBreak = false
				break
			}
			EvaluateStmt(ctx, stmt.Post)
		}
	case *apast.BreakStmt:
		ctx.shouldBreak = true
	case *apast.ReturnStmt:
		returnValues := []reflect.Value{}
		for _, result := range stmt.Results {
			returnValues = append(returnValues, evaluateExpr(ctx, result))
		}
		ctx.returnValues = returnValues
	default:
		panic(fmt.Sprint("Statement eval not implemented: ", reflect.TypeOf(stmt)))
	}
}


func evaluateExpr(ctx *Context, expr apast.Expr) reflect.Value {
	switch expr := expr.(type) {
	case *apast.FuncCallExpr:
		maybeBuiltin := resolveBuiltin(ctx, expr)
		if (maybeBuiltin != nil) {
			return maybeBuiltin()
		}

		funcValue := evaluateExpr(ctx, expr.Func)
		argValues := []reflect.Value{}
		for _, argExpr := range expr.Args {
			argValues = append(argValues, evaluateExpr(ctx, argExpr))
		}
		// TODO: Handle multiple return values.
		return funcValue.Call(argValues)[0]
	case *apast.IdentExpr:
		return ctx.resolveValue(expr.Name)
	case *apast.LiteralExpr:
		return expr.Val
	default:
		panic(fmt.Sprint("Expression eval not implemented: ", reflect.TypeOf(expr)))
	}
}
