package main

import (
	"go/parser"
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
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
	for _, line := range stmt.List {
		switch line := line.(type) {
		case *ast.ExprStmt:
			evaluateExpr(line.X)
		}
	}
}

func evaluateExpr(expr ast.Expr) interface{} {
	switch expr := expr.(type) {
	case *ast.CallExpr:
		funValue := reflect.ValueOf(evaluateExpr(expr.Fun))
		argsValues := []reflect.Value{}
		for _, arg := range expr.Args {
			argsValues = append(argsValues, reflect.ValueOf(evaluateExpr(arg)))
		}
		return funValue.Call(argsValues)
	case *ast.BasicLit:
		switch expr.Kind {
		case token.STRING:
			return parseString(expr.Value)
		default:
			fmt.Print("Unrecognized kind: ", expr.Kind)
			return nil
		}
	case *ast.SelectorExpr:
		if leftSide, ok := expr.X.(*ast.Ident); ok {
			packageName := leftSide.Name
			funcName := expr.Sel.Name
			if packageName == "fmt" && funcName == "Print" {
				return fmt.Print
			}
		}
		return nil
	default:
		fmt.Print("Unexpected token ", reflect.TypeOf(expr), "\n")
	}
	return nil
}

func parseString(codeString string) string {
	strWithoutQuotes := codeString[1:len(codeString) - 1]
	// TODO: Replace with an implementation that properly escapes
	// everything.
	return strings.Replace(strWithoutQuotes, "\\n", "\n", -1)
}