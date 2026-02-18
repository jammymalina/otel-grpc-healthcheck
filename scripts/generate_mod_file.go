package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"text/template"

	"golang.org/x/mod/modfile"
)

type templateData struct {
	OtelVersion      string
	GoVersion        string
	ComponentVersion string
}

const rootDir = ".."

func getGoVersion(goModPath string) (string, error) {
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", goModPath, err)
	}

	f, err := modfile.Parse(goModPath, data, nil)
	if err != nil {
		return "", fmt.Errorf("failed to parse go.mod: %w", err)
	}
	if f.Go == nil {
		return "", errors.New("go version is unset")
	}

	return f.Go.Version, nil
}

func generateFileFromTemplate(tmplFile, targetFile string, data templateData) error {
	tmpl, err := template.New(tmplFile).ParseFiles(path.Join(rootDir, tmplFile))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", tmplFile, err)
	}

	tf, err := os.Create(path.Join(rootDir, targetFile))
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", targetFile, err)
	}
	defer tf.Close()

	err = tmpl.Execute(tf, data)
	if err != nil {
		return fmt.Errorf("failed to execute template into %s: %w", targetFile, err)
	}

	return nil
}

func main() {
	var otelVersion, goModPath, componentVersion string
	flag.StringVar(&otelVersion, "otelversion", "", "OpenTelemetry version")
	flag.StringVar(&goModPath, "otelgomodpath", "", "Path to opentelemetry-collector-contrib go.mod file")
	flag.StringVar(&componentVersion, "componentversion", "", "Component version")

	flag.Parse()

	if otelVersion == "" || componentVersion == "" || goModPath == "" {
		log.Fatal("OpenTelemetry, go.mod path, and Component version must be specified")
	}

	goVersion, err := getGoVersion(goModPath)
	if err != nil {
		log.Fatal("Unable to get go version from opentelemetry-collector-contrib repository")
	}

	data := templateData{
		OtelVersion:      otelVersion,
		GoVersion:        goVersion,
		ComponentVersion: componentVersion,
	}

	filesToGenerate := map[string]string{
		"go.mod.tmpl": "go.mod",
		"readme.tmpl": "README.md",
	}

	for tmpl, target := range filesToGenerate {
		if err := generateFileFromTemplate(tmpl, target, data); err != nil {
			log.Fatalf("Error generating file: %v", err)
		}
	}
}
