package main

import (
	"go/parser"
	"fmt"
	"go/ast"
	"go/token"
	"github.com/alangpierce/apgo/apruntime"
	"github.com/alangpierce/apgo/apcompiler"
	"github.com/alangpierce/apgo/apevaluator"
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
	compiledStmt := apcompiler.CompileStmt(stmt)
	apevaluator.EvaluateStmt(ctx, compiledStmt)
}
