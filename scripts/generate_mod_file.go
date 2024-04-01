package main

import (
	"flag"
	"os"
	"text/template"
)

type templateData struct {
	OtelVersion string
	GoVersion   string
}

func main() {
	var otelVersion, goVersion string
	flag.StringVar(&otelVersion, "otelversion", "", "OpenTelemetry version")
	flag.StringVar(&goVersion, "goversion", "1.20", "Golang version")

	flag.Parse()

	if otelVersion == "" {
		panic("OpenTelemetry version must be specified")
	}

	data := templateData{
		OtelVersion: otelVersion,
		GoVersion:   goVersion,
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
