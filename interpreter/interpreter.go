// Top-level class for managing the interpreter.
package interpreter

import (
	"go/parser"
	"go/token"
	"github.com/alangpierce/apgo/apevaluator"
	"github.com/alangpierce/apgo/apruntime"
	"github.com/alangpierce/apgo/apast"
	"github.com/alangpierce/apgo/apcompiler"
)

type Interpreter struct {
	packages map[string]*apast.Package
	nativePackages map[string]*apruntime.NativePackage
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		packages: make(map[string]*apast.Package),
		nativePackages: make(map[string]*apruntime.NativePackage),
	}
}

// Load and compile the package at the given path.
func (interpreter *Interpreter) LoadPackage(dirPath string) error {
	fset := token.NewFileSet()
	packageAsts, err := parser.ParseDir(fset, dirPath, nil, 0)
	if err != nil {
		return err
	}
	compileCtx := apcompiler.CompileCtx{
		interpreter.nativePackages,
	}
	for name, packageAst := range packageAsts{
		interpreter.packages[name] = apcompiler.CompilePackage(compileCtx, packageAst)
	}
	return nil
}

func (interpreter *Interpreter) LoadNativePackage(pack *apruntime.NativePackage) {
	interpreter.nativePackages[pack.Name] = pack
}

func (interpreter *Interpreter) RunMain() {
	mainPackage := interpreter.packages["main"]
	mainFunc := mainPackage.Funcs["main"]
	ctx := apevaluator.NewContext(mainPackage)
	apevaluator.EvaluateStmt(ctx, mainFunc.Body)
}