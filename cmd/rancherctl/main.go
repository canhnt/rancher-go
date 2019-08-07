package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var VERSION = "dev"

var AppHelpTemplate = `{{.Usage}}
Usage: {{.Name}} {{if .Flags}}[OPTIONS] {{end}}COMMAND [arg...]
Version: {{.Version}}
{{if .Flags}}
Options:
  {{range .Flags}}{{if .Hidden}}{{else}}{{.}}
  {{end}}{{end}}{{end}}
Commands:
  {{range .Commands}}{{.Name}}{{with .Aliases}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
  {{end}}
Run '{{.Name}} COMMAND --help' for more information on a command.
`
var CommandHelpTemplate = `{{.Usage}}
{{if .Description}}{{.Description}}{{end}}
Usage: 
	{{.HelpName}} {{if .Flags}}[OPTIONS] {{end}}{{if ne "None" .ArgsUsage}}{{if ne "" .ArgsUsage}}{{.ArgsUsage}}{{else}}[arg...]{{end}}{{end}}
{{if .Flags}}Options:{{range .Flags}}
	 {{.}}{{end}}{{end}}
`

var SubcommandHelpTemplate = `{{.Usage}}
{{if .Description}}{{.Description}}{{end}}
Usage:
   {{.HelpName}} command{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
Commands:{{range .VisibleCategories}}{{if .Name}}
   {{.Name}}:{{end}}{{range .VisibleCommands}}
     {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}
{{end}}{{if .VisibleFlags}}
Options:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
`

var (
	rancherUrl string
	clusterID  string
	token      string
)

func main() {
	if err := mainErr(); err != nil {
		logrus.Fatal(err)
	}
}

func mainErr() error {
	cli.AppHelpTemplate = AppHelpTemplate
	cli.CommandHelpTemplate = CommandHelpTemplate
	cli.SubcommandHelpTemplate = SubcommandHelpTemplate

	app := cli.NewApp()
	app.Name = "rancherctl"
	app.Usage = "Rancher Project CLI, managing Rancher projects"
	app.Before = func(ctx *cli.Context) error {
		if ctx.GlobalBool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}
		return nil
	}

	app.Version = VERSION
	app.Author = "Canh Ngo"
	app.Email = "canhnt@gmail.com"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Debug logging",
		},
		cli.StringFlag{
			Name:  "rancher-url",
			Usage: "URL of the Rancher server",
		},
		cli.StringFlag{
			Name:  "cluster",
			Usage: "Target cluster ID to manage projects",
		},
		cli.StringFlag{
			Name:  "token",
			Usage: "Security token used to access Rancher APIs",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:        "ls",
			Usage:       "List projects",
			Description: "\nList all projects in the K8s cluster managed by Rancher server",
			ArgsUsage:   "None",
			Action:      defaultAction(projectLs),
		},
		{
			Name:        "get",
			Usage:       "Get project",
			Description: "\nGet project(s) in the K8s cluster managed by Rancher server",
			ArgsUsage:   "ID",
			Action:      defaultAction(projectGet),
		},
		{
			Name:        "apply",
			Usage:       "Create or update multiple projects",
			Description: "\nCreate or update projects defined in the config file to the K8s cluster managed by Rancher server",
			ArgsUsage:   "None",
			Action:      defaultAction(projectApply),
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "filename",
					Usage: "Configuration file containing multiple project information",
				},
			},
		},
		{
			Name:        "delete",
			Usage:       "Remove a project",
			Description: "\nDelete a project from the given k8s cluster managed by Rancher",
			ArgsUsage:   "projectID",
			Action:      defaultAction(projectDelete),
		},
	}

	parsed, err := parseArgs(os.Args)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	logrus.Debugf("Parsed flags: %s", parsed)

	return app.Run(parsed)
}
