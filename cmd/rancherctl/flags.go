package main

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var singleAlphaLetterRegxp = regexp.MustCompile("[a-zA-Z]")

func parseArgs(args []string) ([]string, error) {
	var result []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") && len(arg) > 1 {
			for i, c := range arg[1:] {
				if string(c) == "=" {
					if i < 1 {
						return nil, errors.New("invalid input with '-' and '=' flag")
					}
					result[len(result)-1] = result[len(result)-1] + arg[i+1:]
					break
				} else if singleAlphaLetterRegxp.MatchString(string(c)) {
					result = append(result, "-"+string(c))
				} else {
					return nil, errors.Errorf("invalid input %v in flag", string(c))
				}
			}
		} else {
			result = append(result, arg)
		}
	}
	return result, nil
}

func defaultAction(fn func(ctx *cli.Context) error) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		if ctx.Bool("help") {
			cli.ShowAppHelp(ctx)
			return nil
		}

		rancherUrl = ctx.GlobalString("rancher-url")
		clusterID = ctx.GlobalString("cluster")
		token = ctx.GlobalString("token")
		if err := checkArgs(); err != nil {
			return err
		}

		return fn(ctx)
	}
}

func checkArgs() error {
	if rancherUrl == "" {
		return errors.New("Invalid arguments 'rancher-url'")
	}

	if clusterID == "" {
		return errors.New("Invalid arguments 'cluster'")
	}

	if token == "" {
		return errors.New("Invalid arguments 'token'")
	}

	return nil
}
