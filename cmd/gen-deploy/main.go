package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/remvn/dtalk"
	"github.com/remvn/dtalk/internal/pkg/random"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"text/template"

	"golang.org/x/term"
)

type TemplateData struct {
	AppPort          int
	AppDomain        string
	LiveKitDomain    string
	DnsProviderToken string

	LiveKitClientURL     string
	LiveKitApiKey        string
	LiveKitApiSecret     string
	JwtAccessTokenSecret string
	RedisPassword        string
}

func main() {
	templateDir := "deploy-template"
	entries, err := dtalk.EmbedFS.ReadDir(templateDir)
	if err != nil {
		panic(err)
	}

	templateFiles := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		templateFiles = append(templateFiles, entry.Name())
	}

	outputDir := Prompt{
		name:    "output dir",
		message: "Output directory: ",
	}.readStr()
	if len(outputDir) == 0 {
		panic("output dir is empty")
	}

	templateData := gatherTemplateData()

	if !confirm(outputDir) {
		fmt.Println("aborted.")
		return
	}

	err = os.RemoveAll(outputDir)
	if err != nil {
		panic(fmt.Errorf("Can't remove old output dir: %s: %w", outputDir, err))
	}

	err = os.MkdirAll(outputDir, 0750)
	if err != nil {
		panic(fmt.Errorf("unable to create output dir: %s: %w", outputDir, err))
	}

	fmt.Println()
	fmt.Println("files generated: ")

	for _, name := range templateFiles {
		templateFile := filepath.Join(templateDir, name)
		buffer, err := dtalk.EmbedFS.ReadFile(templateFile)
		if err != nil {
			panic(fmt.Errorf("unable to read file: %s: %w", templateFile, err))
		}
		data := executeTemplate(name, buffer, templateData)
		outputFile := writeFile(name, outputDir, data)
		fmt.Println(outputFile)
	}
}

func confirm(outputDir string) bool {
	answer := Prompt{
		name:    "answer",
		message: fmt.Sprintf("This action will override \"%s\" dir, type \"yes\" to continue: ", outputDir),
	}.readStr()
	return answer == "yes"
}

func gatherTemplateData() TemplateData {
	data := TemplateData{}

	fmt.Println("Leave input prompt empty to use default values")
	data.AppPort = Prompt{
		name:         "app port",
		message:      "App port (8000): ",
		defaultValue: "8000",
	}.readInt()
	data.AppDomain = Prompt{
		name:         "app domain",
		message:      "App domain: (dtalk.yourdomain.com): ",
		defaultValue: "dtalk.yourdomain.com",
	}.readStr()
	data.LiveKitDomain = Prompt{
		name:         "livekit domain",
		message:      "Livekit domain: (livekit.yourdomain.com): ",
		defaultValue: "livekit.yourdomain.com",
	}.readStr()
	data.DnsProviderToken = Prompt{
		name:    "dns provider token",
		message: "DNS provider token (this template is using cloudflare): ",
		hide:    true,
	}.readStr()

	data.LiveKitClientURL = fmt.Sprintf("wss://%s", data.LiveKitDomain)
	data.LiveKitApiKey = random.RandString(15)
	data.LiveKitApiSecret = random.RandString(45)
	data.JwtAccessTokenSecret = random.RandString(45)
	data.RedisPassword = random.RandString(20)
	return data
}

type Prompt struct {
	name         string
	message      string
	defaultValue string
	hide         bool
}

func (p Prompt) readStr() string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(p.message)
	value := ""
	if p.hide {
		arr, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			panic(fmt.Errorf("unable to read %s: %w", p.name, err))
		}
		value = string(arr)
		fmt.Println()
	} else {
		str, err := reader.ReadString('\n')
		if err != nil {
			panic(fmt.Errorf("unable to read %s: %w", p.name, err))
		}
		value = str[:len(str)-1]
	}

	if len(value) == 0 && len(p.defaultValue) > 0 {
		value = p.defaultValue
	}

	return value
}

func (p Prompt) readInt() int {
	str := p.readStr()
	num, err := strconv.Atoi(str)
	if err != nil {
		panic(fmt.Errorf("%s must be an integer: %w", p.name, err))
	}
	return num
}

func executeTemplate(name string, fileBuffer []byte, data TemplateData) []byte {
	t := template.New(name)
	t, err := t.Parse(string(fileBuffer))
	if err != nil {
		panic(fmt.Errorf("unable to parse template %s: %w", name, err))
	}
	buffer := new(bytes.Buffer)
	err = t.Execute(buffer, data)
	if err != nil {
		panic(fmt.Errorf("unable to execute template %s: %w", name, err))
	}
	return buffer.Bytes()
}

func writeFile(name string, dir string, data []byte) string {
	file := filepath.Join(dir, name)
	err := os.WriteFile(file, data, 0640)
	if err != nil {
		panic(fmt.Errorf("unable to write file %s: %w", name, err))
	}
	return file
}
