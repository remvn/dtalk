package main

import (
	"bufio"
	"bytes"
	"dtalk/internal/pkg/random"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"text/template"

	"golang.org/x/term"
)

type AppConfig struct {
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

	outputDir := Prompt{
		name:    "output dir",
		message: "Output directory: ",
	}.readStr()
	config := gatherInput()

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

	for _, name := range files {
		data := executeTemplate(name, templateDir, config)
		writeFile(name, outputDir, data)
	}
}

func confirm(outputDir string) bool {
	answer := Prompt{
		name:    "answer",
		message: fmt.Sprintf("This action will override \"%s\" dir, type \"yes\" to continue: ", outputDir),
	}.readStr()
	return answer == "yes"
}

func gatherInput() AppConfig {
	config := AppConfig{}

	fmt.Println("Leave input prompt empty to use default values")
	config.AppPort = Prompt{
		name:         "app port",
		message:      "App port (8000): ",
		defaultValue: "8000",
	}.readInt()
	config.AppDomain = Prompt{
		name:         "app domain",
		message:      "App domain: (dtalk.yourdomain.com): ",
		defaultValue: "dtalk.yourdomain.com",
	}.readStr()
	config.LiveKitDomain = Prompt{
		name:         "livekit domain",
		message:      "Livekit domain: (livekit.yourdomain.com): ",
		defaultValue: "livekit.yourdomain.com",
	}.readStr()
	config.DnsProviderToken = Prompt{
		name:    "dns provider token",
		message: "DNS provider token (this template is using cloudflare): ",
		hide:    true,
	}.readStr()

	config.LiveKitClientURL = fmt.Sprintf("wss://%s", config.LiveKitDomain)
	config.LiveKitApiKey = random.RandString(15)
	config.LiveKitApiSecret = random.RandString(45)
	config.JwtAccessTokenSecret = random.RandString(45)
	config.RedisPassword = random.RandString(20)
	return config
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
	err := os.WriteFile(filepath.Join(outDir, name), data, 0640)
	if err != nil {
		panic(fmt.Errorf("unable to write file %s: %w", name, err))
	}
}
