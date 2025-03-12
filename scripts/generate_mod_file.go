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

func generateFileFromTemplate(tmplFile, targetFile string, data templateData) {
	tmpl, err := template.New(tmplFile).ParseFiles(tmplFile)
	if err != nil {
		panic(err)
	}

	tf, err := os.Create(targetFile)
	if err != nil {
		panic(err)
	}
	defer tf.Close()

	err = tmpl.Execute(tf, data)
	if err != nil {
		panic(err)
	}
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

	generateFileFromTemplate("go.mod.tmpl", "go.mod", data)
	generateFileFromTemplate("readme.tmpl", "README.md", data)
}
