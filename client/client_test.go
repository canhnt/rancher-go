package client_test

import (
	"os"
	"testing"

	rancher "github.com/canhnt/rancher-go/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	rancherURL = os.Getenv("RANCHER_SERVER")
	token      = os.Getenv("RANCHER_TOKEN")
	clusterID  = os.Getenv("RANCHER_CLUSTER_ID")
	projectID  = os.Getenv("RANCHER_PROJECT_ID")
)

func Test_defaultClient_GetNamespaces(t *testing.T) {
	require.NotEmpty(t, token, "token must not be empty")
	require.NotEmpty(t, clusterID, "clusterID must not be empty")

	client := rancher.NewClient(rancherURL, token)

	namespaces, err := client.GetNamespaces(clusterID)
	require.NoError(t, err)
	assert.NotEmpty(t, namespaces)
	t.Logf("All namespaces: %v\n", namespaces)
}

func Test_defaultClient_GetProjectIDs(t *testing.T) {
	require.NotEmpty(t, token, "token must not be empty")
	require.NotEmpty(t, clusterID, "clusterID must not be empty")

	client := rancher.NewClient(rancherURL, token)

	projects, err := client.GetProjects(clusterID)
	require.NoError(t, err)
	assert.NotEmpty(t, projects)

	t.Logf("Project IDs: %v", projects)
}

func Test_defaultClient__GetNamespacesOfProject(t *testing.T) {
	require.NotEmpty(t, token, "token must not be empty")
	require.NotEmpty(t, clusterID, "clusterID must not be empty")
	require.NotEmpty(t, projectID, "projectID must not be empty")

	client := rancher.NewClient(rancherURL, token)

	namespaces, err := client.GetProjectNamespaces(clusterID, projectID)
	require.NoError(t, err)
	assert.NotEmpty(t, namespaces)

	t.Logf("Namespaces of project '%s': %v\n", projectID, namespaces)
}

func Test_defaultClient__GetProjectGroups(t *testing.T) {
	require.NotEmpty(t, token, "token must not be empty")
	require.NotEmpty(t, projectID, "projectID must not be empty")

	client := rancher.NewClient(rancherURL, token)
	groups, err := client.GetProjectGroups(projectID)
	require.NoError(t, err)
	assert.NotEmpty(t, groups)

	t.Logf("LDAP groups binding to project '%s': %v\n", projectID, groups)
}

func Test_defaultClient__GetClusters(t *testing.T) {
	require.NotEmpty(t, token, "token must not be empty")

	client := rancher.NewClient(rancherURL, token)
	clusters, err := client.GetClusters()
	require.NoError(t, err)
	assert.NotEmpty(t, clusters)

	t.Logf("Rancher clusters: %v", clusters)
}

func Test_defaultClient__GetProjectQuotas(t *testing.T) {
	require.NotEmpty(t, token, "token must not be empty")

	client := rancher.NewClient(rancherURL, token)
	quotas, err := client.GetProjectQuotas(projectID)
	require.NoError(t, err)

	t.Logf("Project quotas: %v", quotas)
}

func Test_defaultClient__CreateProject(t *testing.T) {
	require.NotEmpty(t, token, "token must not be empty")

	client := rancher.NewClient(rancherURL, token)
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

	projectID, err := client.CreateProject(clusterID, project)
	require.NoError(t, err)
	t.Logf("Created project '%s'", projectID)
}

func Test_defaultClient__GetProjectDetail(t *testing.T) {
	require.NotEmpty(t, token, "token must not be empty")

	client := rancher.NewClient(rancherURL, token)
	prj, err := client.GetProjectDetail(projectID)
	require.NoError(t, err)
	t.Logf("Queried project '%+v'", prj)

}

func Test_defaultClient__UpdateProject(t *testing.T) {
	require.NotEmpty(t, token, "token must not be empty")

	client := rancher.NewClient(rancherURL, token)
	prj, err := client.GetProjectDetail(projectID)
	require.NoError(t, err)
	t.Logf("Project before update: %+v", prj)

	prj.PodSecurityPolicyID = "unrestricted"
	// remove 1st member
	prj.Members = prj.Members[1:]
	t.Logf("Members: %+v", prj.Members)

	// add a new member
	prj.Members = append(prj.Members, rancher.Member{
		Type:           rancher.MemberTypeUser,
		PrincipalID:    "openldap_user://cn=ngo500,ou=People,dc=example",
		RoleTemplateID: "project-member",
	})

	err = client.UpdateProject(clusterID, projectID, *prj)
	require.NoError(t, err)

	prj, err = client.GetProjectDetail(projectID)
	require.NoError(t, err)
	t.Logf("Project after update: %+v", prj)

}
