package main

import (
	"bsc/src/commands"
	"bsc/src/exeptions"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

const VERSION_MAJOR = 0
const VERSION_MINOR = 0
const VERSION_PATCH = 2

func main() {
	app := &cli.App{
		Name:        "BSC",
		Description: "BSC is the compiler and package manager for BlockScript. For documentation and exmaples on how to use the language visit: https://block-script.com/docs",
		Version:     fmt.Sprintf("%d.%d.%d", VERSION_MAJOR, VERSION_MINOR, VERSION_PATCH),
		Authors: []*cli.Author{{
			Name:  "H. Schwerdtner",
			Email: "henry.schwerdtner@web.de",
		}},
		Action: commands.DefaultAction,
		Commands: []*cli.Command{
			{
				Name:        "init",
				Aliases:     []string{"i"},
				Usage:       "BSC init [projectName]",
				Description: "Initializes a bsc project and generates the project.json template for you to use.",
				Action:      commands.InitProject,
			},
			{
				Name:        "compile",
				Aliases:     []string{"c"},
				Usage:       "BSC compile [path]",
				Description: "Compiles a project into the specified out folder.",
				Action:      commands.Compile,
			},
			{
				Name:    "versiononly",
				Aliases: []string{"vo"},
				Action: func(c *cli.Context) error {
					fmt.Printf("bsc_v%d.%d.%d\n", VERSION_MAJOR, VERSION_MINOR, VERSION_PATCH)
					return nil
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		if compilerErr, ok := err.(*exeptions.CompilerError); ok {
			fmt.Println(compilerErr.Error())
		} else {
			fmt.Println(err)
		}
	}
}
