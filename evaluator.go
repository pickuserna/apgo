package main

import (
	"go/parser"
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"github.com/alangpierce/apgo/apast"
	"github.com/alangpierce/apgo/apruntime"
	"strconv"
"strings"
)

func main() {
	fset := token.NewFileSet()
	packageAsts, err := parser.ParseDir(fset, "sample", nil, 0)
	if err != nil {
		fmt.Print("Parse error!")
		return
	}
	evaluateMainPackage(packageAsts["main"])
}

func evaluateMainPackage(packageAst *ast.Package) {
	for _, file := range packageAst.Files {
		for _, decl := range file.Decls {
			switch decl := decl.(type) {
			case *ast.FuncDecl:
				if decl.Name.Name == "main" {
					evaluateEmptyFunc(decl.Body)
				}
			}
		}
	}
}

func evaluateEmptyFunc(stmt *ast.BlockStmt) {
	ctx := make(apruntime.Context)
	compiledStmt := compileStmt(stmt)
	evaluateStmt(ctx, compiledStmt)
}

func compileStmt(stmt ast.Stmt) apast.Stmt {
	switch stmt := stmt.(type) {
	//case *ast.BadStmt:
	//	return nil
	//case *ast.DeclStmt:
	//	return nil
	//case *ast.EmptyStmt:
	//	return nil
	//case *ast.LabeledStmt:
	//	return nil
	case *ast.ExprStmt:
		return &apast.ExprStmt{
			compileExpr(stmt.X),
		}
	//case *ast.SendStmt:
	//	return nil
	//case *ast.IncDecStmt:
	//	return nil
	case *ast.AssignStmt:
		lhs := []apast.Expr{}
		rhs := []apast.Expr{}
		for _, lhsExpr := range stmt.Lhs {
			lhs = append(lhs, compileExpr(lhsExpr))
		}
		for _, rhsExpr := range stmt.Rhs {
			rhs = append(rhs, compileExpr(rhsExpr))
		}
		return &apast.AssignStmt{
			lhs,
			rhs,
		}
	//case *ast.GoStmt:
	//	return nil
	//case *ast.DeferStmt:
	//	return nil
	//case *ast.ReturnStmt:
	//	return nil
	//case *ast.BranchStmt:
	//	return nil
	case *ast.BlockStmt:
		stmts := []apast.Stmt{}
		for _, subStmt := range stmt.List {
			stmts = append(stmts, compileStmt(subStmt))
		}
		return &apast.BlockStmt{
			stmts,
		}
	case *ast.IfStmt:
		return nil
	case *ast.CaseClause:
		return nil
	case *ast.SwitchStmt:
		return nil
	case *ast.TypeSwitchStmt:
		return nil
	case *ast.CommClause:
		return nil
	case *ast.SelectStmt:
		return nil
	case *ast.ForStmt:
		return nil
	case *ast.RangeStmt:
		return nil
	default:
		panic(fmt.Sprint("Statement compile not implemented: ", reflect.TypeOf(stmt)))
	}
}

func compileExpr(expr ast.Expr) apast.Expr {
	switch expr := expr.(type) {
	//case *ast.BadExpr:
	//	return nil
	case *ast.Ident:
		return &apast.IdentExpr{
			expr.Name,
		}
	//case *ast.Ellipsis:
	//	return nil
	case *ast.BasicLit:
		return &apast.LiteralExpr{
			reflect.ValueOf(parseLiteral(expr.Value, expr.Kind)),
		}
	//case *ast.FuncLit:
	//	return nil
	//case *ast.CompositeLit:
	//	return nil
	//case *ast.ParenExpr:
	//	return nil
	case *ast.SelectorExpr:
		if leftSide, ok := expr.X.(*ast.Ident); ok {
			packageName := leftSide.Name
			funcName := expr.Sel.Name
			if packageName == "fmt" && funcName == "Print" {
				return &apast.LiteralExpr{
					reflect.ValueOf(fmt.Print),
				}
			}
		}
		return nil
	//case *ast.IndexExpr:
	//	return nil
	//case *ast.SliceExpr:
	//	return nil
	//case *ast.TypeAssertExpr:
	//	return nil
	case *ast.CallExpr:
		compiledArgs := []apast.Expr{}
		for _, arg := range expr.Args {
			compiledArgs = append(compiledArgs, compileExpr(arg))
		}
		return &apast.FuncCallExpr{
			compileExpr(expr.Fun),
			compiledArgs,
		}
	//case *ast.StarExpr:
	//	return nil
	//case *ast.UnaryExpr:
	//	return nil
	case *ast.BinaryExpr:
		return &apast.FuncCallExpr{
			&apast.LiteralExpr{
				apruntime.BinaryOperators[expr.Op],
			},
			[]apast.Expr{compileExpr(expr.X), compileExpr(expr.Y)},
		}
	//case *ast.KeyValueExpr:
	//	return nil
	//
	//case *ast.ArrayType:
	//	return nil
	//case *ast.StructType:
	//	return nil
	//case *ast.FuncType:
	//	return nil
	//case *ast.InterfaceType:
	//	return nil
	//case *ast.MapType:
	//	return nil
	//case *ast.ChanType:
	//	return nil
	default:
		panic(fmt.Sprint("Expression compile not implemented: ", reflect.TypeOf(expr)))
	}
	return nil
}

func evaluateStmt(ctx apruntime.Context, stmt apast.Stmt) {
	switch stmt := stmt.(type) {
	case *apast.ExprStmt:
		evaluateExpr(ctx, stmt.E)
	case *apast.BlockStmt:
		for _, line := range stmt.Stmts {
			evaluateStmt(ctx, line)
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

// parseLiteral takes a primitive literal and returns it as a value.
func parseLiteral(val string, kind token.Token) interface{} {
	switch kind {
	case token.IDENT:
		panic("TODO")
		return nil
	case token.INT:
		// Note that base 0 means that octal and hex literals are also
		// handled.
		result, err := strconv.ParseInt(val, 0, 64)
		if err != nil {
			panic(err)
		}
		return result
	case token.FLOAT:
		panic("TODO")
		return nil
	case token.IMAG:
		panic("TODO")
		return nil
	case token.CHAR:
		panic("TODO")
		return nil
	case token.STRING:
		return parseString(val)
	default:
		fmt.Print("Unrecognized kind: ", kind)
		return nil
	}
}

func parseString(codeString string) string {
	strWithoutQuotes := codeString[1:len(codeString) - 1]
	// TODO: Replace with an implementation that properly escapes
	// everything.
	return strings.Replace(strWithoutQuotes, "\\n", "\n", -1)
}