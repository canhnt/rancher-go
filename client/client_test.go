package client_test

import (
	"os"
	"testing"

	rancher "github.com/canhnt/rancher-go/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const rancherURL = "https://rancher.example.org"

var (
	token     = os.Getenv("RANCHER_TOKEN")
	clusterID = os.Getenv("RANCHER_CLUSTER_ID")
	projectID = os.Getenv("RANCHER_PROJECT_ID")
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
