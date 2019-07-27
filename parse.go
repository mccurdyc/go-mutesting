package mutesting

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
	"path/filepath"

	"golang.org/x/tools/go/packages"
)

// ParseFile parses the content of the given file and returns the corresponding ast.File node and its file set for positional information.
// If a fatal error is encountered the error return argument is not nil.
func ParseFile(file string) (*ast.File, *token.FileSet, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, nil, err
	}

	return ParseSource(data)
}

// ParseSource parses the given source and returns the corresponding ast.File node and its file set for positional information.
// If a fatal error is encountered the error return argument is not nil.
func ParseSource(data interface{}) (*ast.File, *token.FileSet, error) {
	fset := token.NewFileSet()

	src, err := parser.ParseFile(fset, "", data, parser.ParseComments|parser.AllErrors)
	if err != nil {
		return nil, nil, err
	}

	return src, fset, err
}

// ParseAndTypeCheckFile parses and type-checks the given file, and returns everything interesting about the file.
// If a fatal error is encountered the error return argument is not nil.
func ParseAndTypeCheckFile(file string) (*ast.File, *token.FileSet, *types.Package, *types.Info, error) {
	fileAbs, err := filepath.Abs(file)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("Could not absolute the file path of %q: %v", file, err)
	}

	var cfg = &packages.Config{
		Mode: packages.NeedSyntax | packages.NeedFiles,
	}

	pkgs, err := packages.Load(cfg, file)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("Could not load package of file %q: %v", file, err)
	}

	pkgInfo := pkgs[0]

	var src *ast.File
	for _, f := range pkgInfo.Syntax {
		if pkgInfo.Fset.Position(f.Pos()).Filename == fileAbs {
			src = f

			break
		}
	}

	return src, pkgInfo.Fset, pkgInfo.Types, pkgInfo.TypesInfo, nil
}
