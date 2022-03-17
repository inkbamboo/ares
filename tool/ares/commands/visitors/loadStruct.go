package visitors

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"go/ast"
)

type LoadStruct struct{}
type loadVisitor struct {
	Partern string
	Result  []string
}

func (v *loadVisitor) Visit(node ast.Node) ast.Visitor {
	switch node.(type) {
	case *ast.TypeSpec:
		// 遍历全部的接口类型
		typeSpec := node.(*ast.TypeSpec)
		reg := regexp.MustCompile(v.Partern)
		if reg.Match([]byte(typeSpec.Name.Name)) {
			v.Result = append(v.Result, typeSpec.Name.Name)
		}
		return nil
	}
	return v
}

func (v *LoadStruct) GetStructs(pkgPath, partern string) (result []string) {
	ctx := context.TODO()
	pkgs, errs := Load(ctx, pkgPath, os.Environ(), []string{})
	if len(errs) > 0 {
		fmt.Printf("%+v\n", errs)
		return
	}
	for _, pkg := range pkgs {
		for _, item := range pkg.Syntax {
			visitor := loadVisitor{Partern: partern}
			ast.Walk(&visitor, item)
			result = append(result, visitor.Result...)
		}
	}
	return result
}
