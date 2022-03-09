package commands

import (
	"bytes"
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

func modPath(p string) string {
	dir := filepath.Dir(p)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			content, _ := ioutil.ReadFile(filepath.Join(dir, "go.mod"))
			mod := RegexpReplace(`module\s+(?P<name>[\S]+)`, string(content), "$name")
			name := strings.TrimPrefix(filepath.Dir(p), dir)
			name = strings.TrimPrefix(name, string(os.PathSeparator))
			if name == "" {
				return fmt.Sprintf("%s/", mod)
			}
			return fmt.Sprintf("%s/%s/", mod, name)
		}
		parent := filepath.Dir(dir)
		if dir == parent {
			return ""
		}
		dir = parent
	}
}

// RegexpReplace replace regexp
func RegexpReplace(reg, src, temp string) string {
	var result []byte
	pattern := regexp.MustCompile(reg)
	for _, subMatches := range pattern.FindAllStringSubmatchIndex(src, -1) {
		result = pattern.ExpandString(result, temp, src, subMatches)
	}
	return string(result)
}

//go:generate packr2
func createProject() (err error) {
	box := packr.New("all", "./../template/project")
	if err = os.MkdirAll(f.Path, 0755); err != nil {
		return
	}
	for _, name := range box.List() {
		tmpl, _ := box.FindString(name)
		if err = writeOneTmpl(f.Path, name, tmpl); err != nil {
			return
		}
	}
	return
}

//go:generate packr2
func createService() (err error) {
	box := packr.New("service", "./../template/microservices")
	servicePath := filepath.Join(f.Path, "packages", f.ServiceName)
	if err = os.MkdirAll(servicePath, 0755); err != nil {
		return
	}
	for _, name := range box.List() {
		tmpl, _ := box.FindString(name)
		if err = writeOneTmpl(servicePath, name, tmpl); err != nil {
			return
		}
	}
	if err = generate(filepath.Join("./packages", f.ServiceName, "internal/wire/wire.go")); err != nil {
		return
	}
	return
}

func writeOneTmpl(basePath, name, tmpl string) (err error) {
	dir := filepath.Dir(name)
	if dir == "." {
		dir = name
	} else {
		if err = os.MkdirAll(filepath.Join(basePath, dir), 0755); err != nil {
			return
		}
	}
	if strings.HasSuffix(name, ".tmpl") {
		name = strings.TrimSuffix(name, ".tmpl")
	}
	if err = write(filepath.Join(basePath, name), tmpl); err != nil {
		return
	}
	return
}
func formatOneTmpl(tmpl string) ([]byte, error) {
	node, err := parser.ParseExpr(tmpl)
	if err != nil {
		fmt.Println("******111111", tmpl)
		return []byte(tmpl), err
	}
	fset := token.NewFileSet()
	var buf bytes.Buffer
	err = format.Node(&buf, fset, node)
	if err != nil {
		return []byte(tmpl), err
	}
	return buf.Bytes(), err
}
func generate(path string) error {
	cmd := exec.Command("go", "generate", path)
	cmd.Dir = f.Path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func write(path, tpl string) (err error) {

	data, err := parse(tpl)
	if err != nil {
		return
	}
	formatData, _ := formatOneTmpl(data)
	return ioutil.WriteFile(path, formatData, 0644)
}

func parse(s string) (string, error) {
	t, err := template.New("").Parse(s)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err = t.Execute(&buf, f); err != nil {
		return "", err
	}
	return buf.String(), err
}

func buildDir(base string, cmd string, n int) string {
	dirs, err := ioutil.ReadDir(base)
	if err != nil {
		panic(err)
	}
	for _, d := range dirs {
		if d.IsDir() && d.Name() == cmd {
			return path.Join(base, cmd)
		}
	}
	if n <= 1 {
		return base
	}
	return buildDir(filepath.Dir(base), cmd, n-1)
}
