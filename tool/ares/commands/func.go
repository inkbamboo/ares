package commands

import (
	"bytes"
	"github.com/gobuffalo/packr/v2"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

//go:generate packr2
func createProject() (err error) {
	box := packr.New("all", "./../template")
	if err = os.MkdirAll(filepath.Join(f.ProjectPath, f.ServiceName), 0755); err != nil {
		return
	}
	for _, name := range box.List() {
		if name == "go.mod.tmpl" {
			if Exists(filepath.Join(f.ProjectPath, "go.mod")) {
				continue
			}
		}
		tmpl, _ := box.FindString(name)
		if name == "go.mod.tmpl" || name == ".gitignore" {
			if err = writeOneTmpl(f.ProjectPath, name, tmpl); err != nil {
				return
			}
		} else {
			if err = writeOneTmpl(filepath.Join(f.ProjectPath, f.ServiceName), name, tmpl); err != nil {
				return
			}
		}

	}
	if err = generate("internal/wire/wire.go"); err != nil {
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
		//fmt.Println("******111111", tmpl)
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
	cmd.Dir = filepath.Join(f.ProjectPath, f.ServiceName)
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
