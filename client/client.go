package client

import (
	"github.com/golang/glog"
	"gopkg.in/resty.v1"
)

// Brief Rancher entity information
type Entity struct {
	ID   string
	Name string
}

// Client is to query Rancher concepts from Rancher gateway
type Client interface {
	// Returns all clusters in Rancher
	GetClusters() ([]Entity, error)

	// Return all namespaces of the cluster
	GetNamespaces(clusterID string) ([]string, error)

	// Return all projects in the cluster
	GetProjects(clusterID string) ([]Entity, error)

	// Return list of namespaces of the project
	GetProjectNamespaces(clusterID, projectID string) ([]string, error)

	// Return LDAP groups binding to the project
	GetProjectGroups(projectID string) ([]string, error)
}

type defaultClient struct {
	serverURL string
	token     string
}

// NewClient returns a Rancher API client
func NewClient(serverURL, token string) Client {
	return &defaultClient{
		serverURL: serverURL,
		token:     token,
	}
}

func (client defaultClient) GetNamespaces(clusterID string) ([]string, error) {
	resp, err := resty.R().
		SetAuthToken(client.token).
		Get(client.serverURL + "/v3/cluster/" + clusterID + "/namespaces")
	if err != nil {
		glog.Errorf("Failed to query Rancher namespaces: %v", err)
		return nil, err
	}
	body := string(resp.Body()[:])

	return parseValues(body, "data.#.id"), nil
}

func (client defaultClient) GetProjects(clusterID string) ([]Entity, error) {
	resp, err := resty.R().
		SetAuthToken(client.token).
		Get(client.serverURL + "/v3/cluster/" + clusterID + "/projects")
	if err != nil {
		glog.Errorf("Failed to query Rancher projects: %v", err)
		return nil, err
	}
	body := string(resp.Body()[:])
	projects := parseEntities(body, "data")
	return projects, nil
}

func (client defaultClient) GetProjectNamespaces(clusterID, projectID string) ([]string, error) {
	resp, err := resty.R().
		SetAuthToken(client.token).
		Get(client.serverURL + "/v3/cluster/" + clusterID + "/namespaces?projectId=" + projectID)
	if err != nil {
		glog.Errorf("Failed to query Rancher projects: %v", err)
		return nil, err
	}
	body := string(resp.Body()[:])
	return parseValues(body, "data.#.id"), nil
}

// Return LDAP groups binding to the project
func (client defaultClient) GetProjectGroups(projectID string) ([]string, error) {
	resp, err := resty.R().
		SetAuthToken(client.token).
		Get(client.serverURL + "/v3/projects/" + projectID + "/projectroletemplatebindings")
	if err != nil {
		glog.Errorf("Failed to query Rancher project groups: %v", err)
		return nil, err
	}
	body := string(resp.Body()[:])
	groupPrincipalIds := parseValues(body, "data.#.groupPrincipalId")

	var groups []string
	for _, id := range groupPrincipalIds {
		if id != "" {
			g, err := parseGroupFromPrincipalID(id)
			if err != nil {
				glog.V(3).Infof("Invalid principalID '%s': %v", id, err)
			} else if g != "" {
				groups = append(groups, g)
			}
		}
	}
	return groups, nil
}

func (client defaultClient) GetClusters() ([]Entity, error) {
	resp, err := resty.R().
		SetAuthToken(client.token).
		Get(client.serverURL + "/v3/clusters/")
	if err != nil {
		glog.Errorf("Failed to query Rancher clusters: %v", err)
		return nil, err
	}
	body := string(resp.Body()[:])
	clusters := parseEntities(body, "data")
	return clusters, nil
}
