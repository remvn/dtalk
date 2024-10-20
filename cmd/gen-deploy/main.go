package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

type AppConfig struct {
	AppPort              int
	AppDomain            string
	LiveKitDomain        string
	LiveKitApiKey        string
	LiveKitApiSecret     string
	LiveKitClientURL     string
	JwtAccessTokenSecret string
	DnsProviderToken     string
	RedisPassword        string
}

func main() {
	templateDir := "./deploy-template"
	entries, err := os.ReadDir(templateDir)
	if err != nil {
		panic(err)
	}

	files := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		files = append(files, entry.Name())
	}

	config := AppConfig{}

	outputDir := "./deploy-output"
	err = os.MkdirAll(outputDir, 0711)
	if err != nil {
		panic(fmt.Errorf("unable to create output directory: %s: %w", outputDir, err))
	}
	for _, name := range files {
		data := executeTemplate(name, templateDir, config)
		writeFile(name, outputDir, data)
	}
}

func executeTemplate(name string, templateDir string, config AppConfig) []byte {
	t := template.New(name)
	byteArr, err := os.ReadFile(filepath.Join(templateDir, name))
	if err != nil {
		panic(fmt.Errorf("unable to read %s: %w", name, err))
	}
	t, err = t.Parse(string(byteArr))
	if err != nil {
		panic(fmt.Errorf("unable to parse template %s: %w", name, err))
	}
	buffer := new(bytes.Buffer)
	err = t.Execute(buffer, config)
	if err != nil {
		panic(fmt.Errorf("unable to execute template %s: %w", name, err))
	}
	return buffer.Bytes()
}

func writeFile(name string, outDir string, data []byte) {
	err := os.WriteFile(filepath.Join(outDir, name), data, 0644)
	if err != nil {
		panic(fmt.Errorf("unable to write file %s: %w", name, err))
	}
}
