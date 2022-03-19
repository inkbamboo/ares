package visitors

import (
	"bytes"
	"context"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/packages"
	"io/ioutil"
	"os"
)

func GetFile(fileath string, v ast.Visitor) (string, error) {
	file, err := os.Open(fileath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	fileByes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	fSet := token.NewFileSet() // 职位相对于fset
	f, err := parser.ParseFile(fSet, "", string(fileByes), 0)
	if err != nil {
		return "", err
	}
	cMap := ast.NewCommentMap(fSet, f, f.Comments)
	ast.Walk(v, f)
	f.Comments = cMap.Filter(f).Comments()
	var buf bytes.Buffer
	if err = format.Node(&buf, fSet, f); err != nil {
		return "", err
	}
	return string(buf.Bytes()), nil
}

func Load(ctx context.Context, wd string, env []string, patterns []string) ([]*packages.Package, []error) {
	cfg := &packages.Config{
		Context:    ctx,
		Mode:       packages.LoadAllSyntax,
		Dir:        wd,
		Env:        env,
		BuildFlags: []string{},
	}
	escaped := make([]string, len(patterns))
	for i := range patterns {
		escaped[i] = "pattern=" + patterns[i]
	}
	pkgs, err := packages.Load(cfg, escaped...)
	if err != nil {
		return nil, []error{err}
	}
	var errs []error
	for _, p := range pkgs {
		for _, e := range p.Errors {
			errs = append(errs, e)
		}
	}
	if len(errs) > 0 {
		return nil, errs
	}
	return pkgs, nil
}
