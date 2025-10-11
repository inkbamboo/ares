package main

import (
	"fmt"
	"github.com/fatih/color"
)

func main() {

	fmt.Println(color.GreenString("任务列表"))
	fmt.Println(color.Bold + "这是高亮度文本")
}
