package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bitbull-team/bb-aws-connect/cmd/ecs"
	"github.com/bitbull-team/bb-aws-connect/cmd/ssm"

	"github.com/urfave/cli/v2"
)

func main() {
	cmds := []*cli.Command{}
	cmds = append(cmds, ecs.Commands()...)
	cmds = append(cmds, ssm.Commands()...)

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println(
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
			"      ",
		)
		fmt.Fprintf(c.App.Writer, "Bitbull AWS Connect CLI %s\n", c.App.Version)
	}

	cwd, _ := os.Getwd()
	app := &cli.App{
		Name:        "bb-aws-connect",
		Description: "Bitbull AWS Connect CLI",
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
				Value: ".bb-aws-connect.yml",
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
