package main

import (
	"fmt"
	"log"
	"os"

	"app"
	"aws"
	"docker"

	"github.com/urfave/cli/v2"
)

func main() {
	cmds := []*cli.Command{}
	cmds = append(cmds, docker.Commands()...)
	cmds = append(cmds, aws.Commands()...)
	cmds = append(cmds, app.Commands()...)

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println("",
			"		                                               \n",
			"                       ``````                     \n",
			"                        ```````                   \n",
			"                         ````````````:NMMMMMMMd+.`\n",
			"                         `.o:.````````oMMMMMMMMMm`\n",
			"                       ```:Nmh+.```````MMMMMMMMMM`\n",
			"`                     ````.-.``````````MMMMMMMMMM`\n",
			"``               .-/sy+:-.```````````.yMMMMMMMMMM`\n",
			" `            -ohmNMMMMMMmdhyso+++oshNMMMMMMMMMMd`\n",
			" ```       ```dMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMdos`\n",
			" ```  ```````+MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMmoyNM`\n",
			" ```````````oNMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMmmmmd`\n",
			" ````````./dMMMMMMMMMMMMMMMMMMdhhyyso++/::--...`` \n",
			" ``:://ohmMMMMMMMMMMMMMMMMMMMs``                  \n",
			"  `NNMMMMMMMMMMMMMMMMmMMMMMMo`                    \n",
			"  `dMdyyMMm++/::sMM/.`mMMMMo`                     \n",
			"  `sM+`.MM+     .MM.``dMMMo`                      \n",
			"  `/M+`.MM.     `hM.``dMMo`                       \n",
			"   .+-``+/`     `-+` `/+/`\n",
			"      ")
		fmt.Fprintf(c.App.Writer, "Bitbull CLI %s\n", c.App.Version)
	}

	cwd, _ := os.Getwd()
	app := &cli.App{
		Name:        "bb-cli",
		Description: "Bitbull CLI",
		Version:     "VERSION", // this will be overridden during build phase
		Commands:    cmds,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "root",
				Value: cwd,
				Usage: "Change current working directory",
			},
			&cli.StringFlag{
				Name:  "config",
				Value: ".bb-cli.yml",
				Usage: "Config file path",
			},
		},
		Before: func(c *cli.Context) error {
			if c.String("root") != "" {
				os.Chdir(c.String("root"))
				_, err := os.Getwd()
				if err != nil {
					fmt.Println("Cannot change CWD: ", err.Error())
				}
			}
			return nil
		},
		EnableBashCompletion: true,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
