// Top-level class for managing the interpreter.
package interpreter

import (
	"go/parser"
	"go/token"
	"go/ast"
	"github.com/alangpierce/apgo/apevaluator"
	"github.com/alangpierce/apgo/apcompiler"
	"github.com/alangpierce/apgo/apruntime"
)

type Interpreter struct {
	packages map[string]*ast.Package
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		packages: make(map[string]*ast.Package),
	}
}

func (interpreter *Interpreter) LoadPackage(dirPath string) error {
	fset := token.NewFileSet()
	packageAsts, err := parser.ParseDir(fset, dirPath, nil, 0)
	if err != nil {
		return err
	}
	for name, packageAst := range packageAsts{
		interpreter.packages[name] = packageAst
	}
	return nil
}

func (interpreter *Interpreter) RunMain() {
	mainPackage := interpreter.packages["main"]
	for _, file := range mainPackage.Files {
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
