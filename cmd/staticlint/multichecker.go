package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

func main() {
	mychecks := []*analysis.Analyzer{
		//appends.Analyzer,
		//asmdecl.Analyzer,
		//assign.Analyzer,
		//atomic.Analyzer,
		//atomicalign.Analyzer,
		//bools.Analyzer,
		//buildssa.Analyzer,
		//buildtag.Analyzer,
		//cgocall.Analyzer,
		//composite.Analyzer,
		//copylock.Analyzer,
		//ctrlflow.Analyzer,
		//deepequalerrors.Analyzer,
		//defers.Analyzer,
		//directive.Analyzer,
		//errorsas.Analyzer,
		//fieldalignment.Analyzer,
		//findcall.Analyzer,
		//framepointer.Analyzer,
		//httpmux.Analyzer,
		//httpresponse.Analyzer,
		//ifaceassert.Analyzer,
		//inspect.Analyzer,
		//loopclosure.Analyzer,
		//lostcancel.Analyzer,
		//nilfunc.Analyzer,
		//nilness.Analyzer,
		//pkgfact.Analyzer,
		//printf.Analyzer,
		//reflectvaluecompare.Analyzer,
		//shadow.Analyzer,
		//shift.Analyzer,
		//sigchanyzer.Analyzer,
		//slog.Analyzer,
		//sortslice.Analyzer,
		//stdmethods.Analyzer,
		//stdversion.Analyzer,
		//stringintconv.Analyzer,
		//structtag.Analyzer,
		//testinggoroutine.Analyzer,
		//tests.Analyzer,
		//timeformat.Analyzer,
		//unmarshal.Analyzer,
		//unreachable.Analyzer,
		//unsafeptr.Analyzer,
		//unusedresult.Analyzer,
		//unusedwrite.Analyzer,
		//usesgenerics.Analyzer,
		ExitInMainAnalyzer,
	}
	for _, v := range staticcheck.Analyzers {
		mychecks = append(mychecks, v.Analyzer)
	}
	for _, v := range stylecheck.Analyzers {
		mychecks = append(mychecks, v.Analyzer)
	}

	multichecker.Main(
		mychecks...,
	)
}

var ExitInMainAnalyzer = &analysis.Analyzer{
	Name:     "exitinmain",
	Doc:      "check for os.Exit call in func main in package main",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	const MAIN = "main"
	const OS = "os"
	const EXIT = "Exit"
	isMainPkg := func(x *ast.File) bool {
		return x.Name.Name == MAIN
	}

	isMainFunc := func(x *ast.FuncDecl) bool {
		return x.Name.Name == MAIN
	}

	isOsExit := func(x *ast.SelectorExpr, isMain bool) bool {
		if !isMain || x.X == nil {
			return false
		}
		ident, ok := x.X.(*ast.Ident)
		if !ok {
			return false
		}
		if ident.Name == OS && x.Sel.Name == EXIT {
			pass.Reportf(ident.NamePos, "os.Exit called in main func in main package")
			return true
		}
		return false
	}

	i := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.File)(nil),
		(*ast.FuncDecl)(nil),
		(*ast.SelectorExpr)(nil),
	}
	mainInspecting := false
	i.Preorder(nodeFilter, func(n ast.Node) {
		switch x := n.(type) {
		case *ast.File:
			if !isMainPkg(x) {
				return
			}
		case *ast.FuncDecl:
			f := isMainFunc(x)
			if mainInspecting && !f {
				mainInspecting = false
				return
			}
			mainInspecting = f
		case *ast.SelectorExpr:
			if isOsExit(x, mainInspecting) {
				return
			}
		}
	})
	return nil, nil
}
