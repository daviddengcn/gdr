package gdrf

import (
	"bytes"
	"fmt"
	"github.com/daviddengcn/go-villa"
	"go/ast"
	"go/build"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"os"
	"strconv"
)

const (
	pt_VAR = iota
	pt_TYPE
)

type Placeholder struct {
	tp   int
	name string
}

// findPlaceholderInFile finds a placeholder in a ast.File
func findPlaceholderInFile(fs *token.FileSet, f *ast.File) *Placeholder {
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			switch d.Tok {
			case token.TYPE:
				for i := range d.Specs {
					spec := d.Specs[i].(*ast.TypeSpec)
					name := spec.Name.Name
					if ast.IsExported(name) {
						// Placeholder found
						return &Placeholder{tp: pt_TYPE, name: name}
					}
				}
			case token.CONST, token.VAR:
				for i := range d.Specs {
					spec := d.Specs[i].(*ast.ValueSpec)
					for _, ident := range spec.Names {
						name := ident.Name
						if ast.IsExported(name) {
							// Placeholder found
							return &Placeholder{tp: pt_VAR, name: name}
						}
					}
				} // for i
			}
		case *ast.FuncDecl:
			if d.Recv != nil {
				// ignore methods
				continue
			}

			name := d.Name.Name
			if ast.IsExported(name) {
				// Placeholder found
				return &Placeholder{tp: pt_VAR, name: name}
			}
		}
	}

	return nil
}


// findPlaceholder finds a placeholder with a name and path
func findPlaceholder(name, path string) *Placeholder {
	path, err := strconv.Unquote(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	pkg, err := build.Import(path, ".", 0)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	files := villa.NewStrSet(pkg.GoFiles...)

	fs := token.NewFileSet()
	pkgs, err := parser.ParseDir(fs, pkg.Dir, func(fi os.FileInfo) bool {
		if fi.IsDir() {
			return false
		}

		return files.In(fi.Name())
	}, 0)
	if err != nil {
		return nil
	}

	p, ok := pkgs[pkg.Name]
	if !ok {
		fmt.Println(pkg.Name, "not found in", pkgs)
		return nil
	}

	if name == "" {
		name = p.Name
	}
	for _, f := range p.Files {
		ph := findPlaceholderInFile(fs, f)
		if ph != nil {
			ph.name = name + "." + ph.name
			return ph
		}
	}
	return nil
}

// FilterFile parses the go file of inFn and output a file with placeholders
// appended.
//
// If src != nil, FilterFile parses the source from src and the filename is
// only used when recording position information. The type of the argument
// for the src parameter must be string, []byte, or io.Reader.
// If src == nil, FilterFile parses the file specified by filename.
func FilterFile(inFn villa.Path, src interface{}, out io.Writer) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, inFn.S(), src, 0)
	if err != nil {
		return err
	} // if

	var phVars, phTypes villa.StrSet

	for _, imp := range f.Imports {
		name := ""
		if imp.Name != nil {
			name = imp.Name.Name
		}

		if name == "." || name == "_" {
			// no need to hold
			continue
		}
		ph := findPlaceholder(name, imp.Path.Value)
		if ph != nil {
			switch ph.tp {
			case pt_VAR:
				phVars.Put(ph.name)

			case pt_TYPE:
				phTypes.Put(ph.name)
			}
		}
	}

	var initFunc bytes.Buffer
	initFunc.WriteString("\n")
	if len(phVars) > 0 {
		initFunc.WriteString("func _() {\n")
		for v := range phVars {
			initFunc.WriteString("\t_ = " + v + "\n")
		}
		initFunc.WriteString("}\n")
	}

	if len(phTypes) > 0 {
		initFunc.WriteString("type (\n")
		for t := range phTypes {
			initFunc.WriteString("\t_ " + t + "\n")
		}
		initFunc.WriteString(")\n")
	}
	(&printer.Config{Mode: printer.RawFormat, Tabwidth: 4}).Fprint(out, fset, f)

	out.Write(initFunc.Bytes())
	return nil
}
