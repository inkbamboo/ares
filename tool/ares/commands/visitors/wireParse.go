package visitors

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
)

type WireVisitor struct {
}

func (v *WireVisitor) Visit(node ast.Node) ast.Visitor {
	switch node.(type) {
	case *ast.GenDecl:
		genDecl := node.(*ast.GenDecl)
		// 查找有没有import context包
		// Notice：没有考虑没有import任何包的状况
		if genDecl.Tok == token.IMPORT {
			v.addImport(genDecl)
		} else if genDecl.Tok == token.VAR {
			return v
		}
		return nil
	case *ast.ValueSpec:
		valDecl := node.(*ast.ValueSpec)
		if valDecl.Names[0].Name == "svcProvider" {
			return v
		}
		return nil
	case *ast.CallExpr:
		callDecl := node.(*ast.CallExpr)
		v.addService(callDecl)
		return nil
	}

	return v
}

// addImport 引入context包
func (v *WireVisitor) addImport(genDecl *ast.GenDecl) {
	// 是否已经import
	hasImported := false
	for _, v := range genDecl.Specs {
		imptSpec := v.(*ast.ImportSpec)
		// 若是已经包含"context"
		fmt.Println(imptSpec.Path.Value)
		if imptSpec.Path.Value == strconv.Quote("context") {
			hasImported = true
		}
	}
	// 若是没有import context，则import
	if !hasImported {
		genDecl.Specs = append(genDecl.Specs, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: strconv.Quote("context"),
			},
		})
	}
}

// addContext 添加context参数
func (v *WireVisitor) addService(callSpec *ast.CallExpr) {
	// 接口方法不为空时，遍历接口方法
	//for _, argItem := range callSpec.Args {
	//arg := argItem.(*ast.SelectorExpr)
	//if arg.Sel.Name != "NewBoxService" {
	//
	//}
	//}
	newArg := &ast.SelectorExpr{
		X:   ast.NewIdent("service"),
		Sel: ast.NewIdent("NewGuessService"),
	}
	callSpec.Args = append(callSpec.Args, newArg)
}
