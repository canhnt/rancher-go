package main

import (
	"errors"
	"fmt"

	"github.com/Sirupsen/logrus"
	rancher "github.com/canhnt/rancher-go/client"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

func projectApply(ctx *cli.Context) error {
	configFile := ctx.String("filename")
	if configFile == "" {
		return errors.New("config file argument not found")
	}

	client := rancher.NewClient(rancherUrl, token)

	projectList, err := rancher.ReadProjects(configFile)
	if err != nil {
		return err
	}

	success := true
	for _, prj := range projectList.Projects {
		if prj.ID != "" {
			// existed project, update it
			logrus.Infof("Updating project ID='%s', Name='%s'", prj.ID, prj.Name)
			err = client.UpdateProject(clusterID, prj.ID, prj)
			if err != nil {
				logrus.Errorf("Failed to update project ID='%s', name='%s': %v", prj.ID, prj.Name, err)
			}
		} else {
			// new project
			logrus.Infof("Creating project Name='%s'", prj.Name)
			prjID, err := client.CreateProject(clusterID, prj)
			if err != nil {
				logrus.Errorf("Created project '%s' failed: %v", prj.Name, err)
				success = false
			} else {
				logrus.Infof("Created project name='%s', ID='%s'", prj.Name, prjID)
			}
		}
	}
	if !success {
		return errors.New("created projects failed")
	}
	return nil
}

func projectLs(ctx *cli.Context) error {
	client := rancher.NewClient(rancherUrl, token)
	projects, err := client.GetProjects(clusterID)
	if err != nil {
		return err
	}

	fmt.Printf("Projects in cluster '%s'\n", clusterID)
	fmt.Println("ID \t\t\t Name")
	for _, prj := range projects {
		fmt.Printf("%s \t %s\n", prj.ID, prj.Name)
	}
	return nil
}

func projectGet(ctx *cli.Context) error {
	client := rancher.NewClient(rancherUrl, token)
	args := ctx.Args()
	if len(args) == 0 {
		logrus.Debug("Query all projects in cluster")
		projectEntities, err := client.GetProjects(clusterID)
		if err != nil {
			return err
		}
		var projectList rancher.ProjectList
		for _, e := range projectEntities {
			proj, err := client.GetProjectDetail(e.ID)
			if err != nil {
				logrus.Error(err)
			} else {
				projectList.Projects = append(projectList.Projects, *proj)
			}
		}

		err = printYAML(projectList)
		if err != nil {
			return err
		}
	} else {
		project, err := client.GetProjectDetail(args[0])
		if err != nil {
			return err
		}
		err = printYAML(project)
		if err != nil {
			return err
		}
	}
	return nil
}

func projectDelete(ctx *cli.Context) error {
	fmt.Printf("Args delete project: %+v", ctx.Args())
	panic("Implement me")
}

func printYAML(in interface{}) error {
	output, err := yaml.Marshal(in)
	if err != nil {
		return err
	}
	fmt.Print(string(output[:]))
	return nil
}
