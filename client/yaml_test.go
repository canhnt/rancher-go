package client_test

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	rancher "github.com/canhnt/rancher-go/client"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func Test_SaveYAML(t *testing.T) {
	project := rancher.Project{
		Name:        "canh-test-project",
		Description: "Test project created via rancher-go",
		Members: []rancher.Member{
			{
				PrincipalID:    "openldap_group://cn=developers,ou=Groups,dc=example",
				Type:           rancher.MemberTypeGroup,
				RoleTemplateID: "rt-pqglq",
			},
			{
				PrincipalID:    "openldap_group://cn=reviewers,ou=Groups,dc=example",
				Type:           rancher.MemberTypeGroup,
				RoleTemplateID: "project-member",
			},
			{
				PrincipalID:    "openldap_user://cn=ngo500,ou=People,dc=example",
				Type:           rancher.MemberTypeUser,
				RoleTemplateID: "project-member",
			},
		},
		PodSecurityPolicyID: "tcloud",
		ResourceQuotas: rancher.ProjectQuotas{
			Project: map[string]string{
				"limitsCpu":       "2500m",
				"requestsCpu":     "1200m",
				"limitsMemory":    "4096Mi",
				"requestsMemory":  "2048Mi",
				"requestsStorage": "80Gi",
			},
			Namespace: map[string]string{
				"limitsCpu":       "500m",
				"limitsMemory":    "1Gi",
				"requestsCpu":     "200m",
				"requestsMemory":  "512Mi",
				"requestsStorage": "20Gi",
			},
		},
	}

	d, err := yaml.Marshal(&project)
	require.NoError(t, err)

	t.Logf("Marshal:\n%s", string(d[:]))
}

func Test_ReadYAMLProject(t *testing.T) {
	data := `
name: "demo-project"
description: "A sample project storing in YAML file"
podSecurityPolicyId: tcloud
members:
- type: Group
  principalId: openldap_group://cn=developers,ou=Groups,dc=example
  roleTemplateId: project-member
- type: Group
  principalId: openldap_group://cn=testers,ou=Groups,dc=example
  roleTemplateId: rt-123
- type: User
  principalId: openldap_user://cn=canh,ou=Groups,dc=example
  roleTemplateId: project-owner
projectQuotas:
  project:
    limitsCpu: '2500m'
    limitsMemory: '4096Mi'
    requestsCpu: 1200m
    requestsMemory: 2048Mi
    requestsStorage: 80Gi
  namespace:
    limitsCpu: 500m
    limitsMemory: 1Gi
    requestsCpu: 200m
    requestsMemory: 512Mi
    requestsStorage: 20Gi
`
	p := rancher.Project{}
	err := yaml.Unmarshal([]byte(data), &p)
	require.NoError(t, err)
	t.Logf("%+v\n\n", p)
}

func Test_ReadYAMLProjectList(t *testing.T) {
	data := []byte(`
projects:
  - name: "demo-project1"
    description: "1st project storing in YAML file"
    podSecurityPolicyId: tcloud
    members:
    - type: Group
      principalId: openldap_group://cn=developers,ou=Groups,dc=example
      roleTemplateId: project-member
    - type: Group
      principalId: openldap_group://cn=testers,ou=Groups,dc=example
      roleTemplateId: rt-123
    - type: User
      principalId: openldap_user://cn=canh,ou=Groups,dc=example
      roleTemplateId: project-owner
    projectQuotas:
      project:
        limitsCpu: '2500m'
        limitsMemory: '4096Mi'
        requestsCpu: 1200m
        requestsMemory: 2048Mi
        requestsStorage: 80Gi
      namespace:
        limitsCpu: 500m
        limitsMemory: 1Gi
        requestsCpu: 200m
        requestsMemory: 512Mi
        requestsStorage: 20Gi
  - name: "demo-project2"
    description: "2nd project storing in YAML file"
    podSecurityPolicyId: tcloud
    members:
    - type: Group
      principalId: openldap_group://cn=developers,ou=Groups,dc=example
      roleTemplateId: project-member
    - type: Group
      principalId: openldap_group://cn=testers,ou=Groups,dc=example
      roleTemplateId: rt-123
    - type: User
      principalId: openldap_user://cn=canh,ou=Groups,dc=example
      roleTemplateId: project-owner
    projectQuotas:
      project:
        limitsCpu: '2500m'
        limitsMemory: '4096Mi'
        requestsCpu: 1200m
        requestsMemory: 2048Mi
        requestsStorage: 80Gi
      namespace:
        limitsCpu: 500m
        limitsMemory: 1Gi
        requestsCpu: 200m
        requestsMemory: 512Mi
        requestsStorage: 20Gi
`)
	projects := rancher.ProjectList{}

	err := yaml.Unmarshal([]byte(data), &projects)
	require.NoError(t, err)
	t.Logf("%+v\n\n", projects)
}

func Test_WriteAMLProjectList(t *testing.T) {
	var projects = rancher.ProjectList{
		Projects: []rancher.Project{
			{
				Name:        "project-1",
				Description: "1st project",
			},
			{
				Name:        "project-2",
				Description: "2nd project",
			},
		},
	}

	d, err := yaml.Marshal(&projects)
	require.NoError(t, err)
	fmt.Printf("Marshal:\n%s", string(d[:]))
}

func Test_ReadProjects(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	testDir := filepath.Dir(filename)
	fmt.Println("Current test dir: " + testDir)

	projects, err := rancher.ReadProjects(testDir + "/project.yaml")
	require.NoError(t, err)
	t.Logf("%+v\n\n", projects)
}
