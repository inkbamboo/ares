package commands

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
	"path/filepath"
)

func NewAction() cli.ActionFunc {
	return func(ctx *cli.Context) (err error) {
		installToolList()
		err = newProject(ctx)
		//checkGoMode(f.Path)
		return
	}
}

func checkGoMode(dir string) {
	if !Exists(filepath.Join(dir, "go.mod")) {
		RunCmd("init", dir, "go", []string{"mod", "init", filepath.Base(dir)})
		RunCmd("tidy", dir, "go", []string{"mod", "tidy"})
	}
	return
}

func newProject(ctx *cli.Context) (err error) {
	//服务名称
	f.ServiceName = ctx.Args().Get(0)
	if f.ProjectPath != "" {
		if f.ProjectPath, err = filepath.Abs(f.ProjectPath); err != nil {
			return
		}
		f.ProjectName = filepath.Base(f.ProjectPath)
	} else {
		f.ProjectPath, _ = os.Getwd()
	}
	//取当前文件夹为工程名称
	f.ProjectName = filepath.Base(f.ProjectPath)
	f.MysqlList = []DbClient{
		{
			ConfigName: "master",
			ClientName: "MasterClient",
		}, {
			ConfigName: "read",
			ClientName: "ReadClient",
		},
	}
	f.RedisList = []DbClient{{
		ConfigName: "master",
		ClientName: "MasterClient",
	}}
	// creata a project
	if err = createProject(); err != nil {
		return
	}
	fmt.Printf("Project: %s\n", f.ProjectName)
	fmt.Printf("Directory: %s\n", f.ProjectPath)
	fmt.Println("项目创建成功.")
	return
}

func installToolList() {
	for _, t := range toolList() {
		if !t.installed() || t.needUpdated() {
			t.install()
		}
	}
}

func BuildAction(c *cli.Context) error {
	base, err := os.Getwd()
	if err != nil {
		return err
	}
	args := append([]string{"build", "-i"}, c.Args().Slice()...)
	cmd := exec.Command("go", args...)
	cmd.Dir = buildDir(base, "cmd", 5)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Printf("directory: %s\n", cmd.Dir)
	fmt.Printf("ares: %s\n", Version)
	if err := cmd.Run(); err != nil {
		return err
	}
	fmt.Println("build success.")
	return nil
}

func RunAction(c *cli.Context) error {
	base, err := os.Getwd()
	if err != nil {
		return err
	}
	dir := buildDir(base, "cmd", 5)
	//conf := path.Join(filepath.Dir(dir), "configs")
	//args := append([]string{"run", "main.go", "-conf", conf}, c.Args().Slice()...)
	args := append([]string{"run", "main.go"}, c.Args().Slice()...)
	cmd := exec.Command("go", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
