package client

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// ReadProjects reads the YAML file containing list of projects
func ReadProjects(yamlFile string) (*ProjectList, error) {
	data, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return nil, err
	}

	projects := ProjectList{}

	err = yaml.Unmarshal([]byte(data), &projects)
	if err != nil {
		return nil, err
	}

	return &projects, nil
}
