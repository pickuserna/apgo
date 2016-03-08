package apcompiler

import (
	"go/ast"
	"github.com/alangpierce/apgo/apast"
	"go/token"
	"github.com/alangpierce/apgo/apruntime"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func CompileStmt(stmt ast.Stmt) apast.Stmt {
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
		if stmt.Tok == token.DEFINE || stmt.Tok == token.ASSIGN {
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
		} else {
			if len(stmt.Lhs) != 1 || len(stmt.Rhs) != 1 {
				panic("Unexpected multiple assign")
			}
			// TODO: We should only evaluate the left side once,
			// e.g. array index values.
			compiledLhs := compileExpr(stmt.Lhs[0])
			return &apast.AssignStmt{
				[]apast.Expr{compiledLhs},
				[]apast.Expr{
					&apast.FuncCallExpr{
						&apast.LiteralExpr{
							apruntime.AssignBinaryOperators[stmt.Tok],
						},
						[]apast.Expr{
							compiledLhs,
							compileExpr(stmt.Rhs[0]),
						},
					},
				},
			}
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
			stmts = append(stmts, CompileStmt(subStmt))
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