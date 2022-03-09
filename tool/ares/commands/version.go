package commands

import (
	"bytes"
	"runtime"
	"text/template"
	"time"
)

var (
	Version   = "v0.0.1"
	BuildTime = time.Now().Format("2006-01-02 15:04:05")
)

// VersionOptions include version
type VersionOptions struct {
	Version   string
	BuildTime string
	GoVersion string
	Os        string
	Arch      string
}

var versionTemplate = `Version: {{.Version}} | Go version: {{.GoVersion}} | BuildTime: {{.BuildTime}} | OS/Arch: {{.Os}}/{{.Arch}}`

func GetVersion() string {
	var doc bytes.Buffer
	vo := VersionOptions{
		Version:   Version,
		BuildTime: BuildTime,
		GoVersion: runtime.Version(),
		Os:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
	tmpl, _ := template.New("version").Parse(versionTemplate)
	_ = tmpl.Execute(&doc, vo)
	return doc.String()
}
