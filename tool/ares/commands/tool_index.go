package commands

import "time"

var toolIndexs = []*Tool{
	{
		Name:      "ares",
		Alias:     "ares",
		BuildTime: time.Date(2020, 3, 31, 0, 0, 0, 0, time.Local),
		Install:   "go get -u github.com/inkbamboo/ares/tool/ares@" + Version,
		Summary:   "工具集",
		Platform:  []string{"darwin", "linux", "windows"},
		Author:    "ares",
	},
	{
		Name:      "wire",
		Alias:     "wire",
		BuildTime: time.Date(2020, 3, 31, 0, 0, 0, 0, time.Local),
		Install:   "go get -u github.com/google/wire/cmd/wire",
		Platform:  []string{"darwin", "linux", "windows"},
		Author:    "google",
	},
	{
		Name:      "packr2",
		Alias:     "packr2",
		BuildTime: time.Date(2020, 3, 31, 0, 0, 0, 0, time.Local),
		Install:   "go get -u github.com/gobuffalo/packr/v2/packr2",
		Platform:  []string{"darwin", "linux", "windows"},
		Author:    "gobuffalo",
	},
	{
		Name:      "swag",
		Alias:     "swag",
		BuildTime: time.Date(2020, 3, 31, 0, 0, 0, 0, time.Local),
		Install:   "go get -u github.com/swaggo/swag/cmd/swag",
		Platform:  []string{"darwin", "linux", "windows"},
		Author:    "swaggo",
	},
	//{
	//	Name:      "kitex",
	//	Alias:     "kitex",
	//	BuildTime: time.Date(2022, 2, 21, 0, 0, 0, 0, time.Local),
	//	Install:   "go install github.com/cloudwego/kitex/tool/cmd/kitex@v0.1.4",
	//	Platform:  []string{"darwin", "linux", "windows"},
	//	Author:    "kitex",
	//},
	//{
	//	Name:      "thriftgo",
	//	Alias:     "thriftgo",
	//	BuildTime: time.Date(2022, 2, 21, 0, 0, 0, 0, time.Local),
	//	Install:   "go install github.com/cloudwego/thriftgo@v0.1.4",
	//	Platform:  []string{"darwin", "linux", "windows"},
	//	Author:    "thriftgo",
	//},
}
