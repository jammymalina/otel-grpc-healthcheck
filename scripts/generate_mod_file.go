package main

import (
	"flag"
	"fmt"
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
		fmt.Println("OpenTelemetry version must be specified")
		os.Exit(2)
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
	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}
