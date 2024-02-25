package commands

import (
	"bsc/src/compiler"
	"bsc/src/parser"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

func LoadProjectSettings() (*ProjectSettings, error) {
	file, err := os.Open("project.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var settings ProjectSettings
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&settings)
	if err != nil {
		return nil, err
	}

	return &settings, nil
}

func InitProject(ctx *cli.Context) error {
	projectName := ctx.Args().First()
	if projectName == "" {
		return fmt.Errorf("please provide a project name")
	}

	fmt.Println("Initializing Project:", projectName)

	err := os.Mkdir(projectName, 0755)
	if err != nil {
		return err
	}

	filePath := fmt.Sprintf("%s/project.json", projectName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	defaultSettings := &ProjectSettings{
		Name:         projectName,
		Version:      "0.1.0",
		BuildDir:     "build",
		ArtifactName: "out",
		Entry:        "src/main.bs",
		Deps:         make([]string, 0),
		Commands: map[string]string{
			"build": "bsc compile src/main.bs",
			"run":   "bsc compile src/main.bs && ./build/out",
		},
	}
	err = os.Mkdir(projectName+"/src", 0755)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	encoder.SetEscapeHTML(false)
	err = encoder.Encode(defaultSettings)
	if err != nil {
		return err
	}

	file, err = os.Create(projectName + "/src/main.bs")
	if err != nil {
		return err
	}
	defer file.Close()
	_, writeErr := file.WriteString("exit(0);")
	return writeErr
}

func Compile(ctx *cli.Context) error {
	settings, err := LoadProjectSettings()
	var buildDirName string
	var artifactName string
	if err == nil {
		buildDirName = "build"
		artifactName = "out"
	} else {
		buildDirName = settings.BuildDir
		artifactName = settings.ArtifactName
	}
	os.RemoveAll(buildDirName)
	os.Mkdir(buildDirName, 0755)
	entryPoint := ctx.Args().First()
	if entryPoint == "" && err != nil {
		return errors.New("no src file declared in command or in the project.json")
	}
	file, entryErr := os.Open(entryPoint)
	if entryErr != nil {
		return errors.New("file does not exist")
	}
	parser := parser.NewNParser()
	ast, parsingError := parser.Parse(file.Name(), file)
	if parsingError != nil {
		return parsingError
	}
	compiler := compiler.NewBQCCompiler(ast)
	return compiler.Compile(buildDirName, artifactName)
}

func DefaultAction(ctx *cli.Context) error {
	if ctx.Args().Present() {
		command := ctx.Args().First()
		settings, err := LoadProjectSettings()
		if err != nil {
			return errors.New("the provided command cannot be found in the project.json")
		}
		toExecute, exists := settings.Commands[command]
		if !exists {
			return errors.New("the provided command cannot be found in the project.json")
		}
		err = exec.Command("sh", "-c", toExecute).Run()
		if err != nil {
			return errors.New(err.Error())
		}
	} else {
		return errors.New("no command provided")
	}
	return nil
}
