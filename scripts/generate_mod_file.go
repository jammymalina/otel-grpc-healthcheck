package main

import (
	"flag"
	"os"
	"text/template"
)

type templateData struct {
	OtelVersion      string
	GoVersion        string
	ComponentVersion string
}

func main() {
	var otelVersion, goVersion, componentVersion string
	flag.StringVar(&otelVersion, "otelversion", "", "OpenTelemetry version")
	flag.StringVar(&goVersion, "goversion", "1.20", "Golang version")
	flag.StringVar(&componentVersion, "componentversion", "", "Component version")

	flag.Parse()

	if otelVersion == "" || componentVersion == "" {
		panic("OpenTelemetry and Component version must be specified")
	}

	data := templateData{
		OtelVersion:      otelVersion,
		GoVersion:        goVersion,
		ComponentVersion: componentVersion,
	}

	tmplFile := "go.mod.tmpl"
	tmpl, err := template.New(tmplFile).ParseFiles(tmplFile)
	if err != nil {
		panic(err)
	}

	goModFile, err := os.Create("go.mod")
	if err != nil {
		panic(err)
	}
	defer goModFile.Close()

	err = tmpl.Execute(goModFile, data)
	if err != nil {
		panic(err)
	}
}
