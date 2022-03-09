package commands

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Tool is kratos tool.
type Tool struct {
	Name      string    `json:"name"`
	Alias     string    `json:"alias"`
	BuildTime time.Time `json:"build_time"`
	Install   string    `json:"install"`
	Dir       string    `json:"dir"`
	Summary   string    `json:"summary"`
	Platform  []string  `json:"platform"`
	Author    string    `json:"author"`
	URL       string    `json:"url"`
}

func toolList() (tools []*Tool) {
	return toolIndexs
}

func (t Tool) needUpdated() bool {
	if !t.supportOS() || t.Install == "" {
		return false
	}
	if f, err := os.Stat(t.toolPath()); err == nil {
		if t.BuildTime.After(f.ModTime()) {
			return true
		}
	}
	return false
}

func (t Tool) toolPath() string {
	name := t.Alias
	if name == "" {
		name = t.Name
	}
	gobin := Getenv("GOBIN")
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	if gobin != "" {
		return filepath.Join(gobin, name)
	}
	return filepath.Join(GetGOPATH(), "bin", name)
}

func (t Tool) installed() bool {
	_, err := os.Stat(t.toolPath())
	return err == nil
}

func (t Tool) supportOS() bool {
	for _, p := range t.Platform {
		if strings.ToLower(p) == runtime.GOOS {
			return true
		}
	}
	return false
}

func (t Tool) install() {
	if t.Install == "" {
		fmt.Fprintf(os.Stderr, color.RedString("%s: 安装失败\n", t.Name))
		return
	}
	fmt.Println(t.Install)
	cmds := strings.Split(t.Install, " ")
	if len(cmds) > 0 {
		if err := RunCmd(t.Name, path.Dir(t.toolPath()), cmds[0], cmds[1:]); err == nil {
			color.Green("%s: 安装成功!", t.Name)
		}
	}
}

func (t Tool) updated() bool {
	if !t.supportOS() || t.Install == "" {
		return false
	}
	if f, err := os.Stat(t.toolPath()); err == nil {
		if t.BuildTime.After(f.ModTime()) {
			return true
		}
	}
	return false
}
