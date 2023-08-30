package commands

import (
	"bsc/src/compiler"
	"bsc/src/lexer"
	"bsc/src/parser"
	"encoding/json"
	"errors"
	"fmt"
	"os"

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
		return fmt.Errorf("Please provide a project name")
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
		Commands: map[string]interface{}{
			"build": "bsc compile src/main.bs",
			"run":   "bsc compile src/main.bs && ./build/out",
		},
	}

	data, err := json.MarshalIndent(defaultSettings, "", "    ")
	if err != nil {
		return err
	}

	_, writeErr := file.Write(data)
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
		return errors.New("No src file declared in command or in the project.json.")
	}
	file, entryErr := os.Open(entryPoint)
	if entryErr != nil {
		return errors.New("File does not exist.")
	}
	tokenizer := lexer.NewTokenizer(file)
	parser := parser.NewParser(tokenizer)
	ast, parsingError := parser.Parse()
	if parsingError != nil {
		return parsingError
	}
	compiler := compiler.NewNASMElf64Compiler(ast)
	return compiler.Compile(buildDirName, artifactName)
}
