package client

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"gopkg.in/resty.v1"
)

// Brief Rancher entity information
type Entity struct {
	ID   string
	Name string
}

// Reader is to query Rancher concepts from Rancher gateway
type Reader interface {
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

	// Return quotas set in the project, e.g. CPU, memory, storage, etc.
	GetProjectQuotas(projectID string) (*ProjectQuotas, error)

	// Return members of project
	GetProjectMembers(projectID string) ([]Member, error)

	GetProjectDetail(projectID string) (*Project, error)
}

type Writer interface {
	CreateProject(clusterID string, project Project) (string, error)

	UpdateProject(clusterID, projectID string, project Project) error

	AddProjectMember(projectID string, member Member) (err error)
}

type Client interface {
	Reader
	Writer
}

type defaultClient struct {
	serverURL string
	token     string
}

func (client defaultClient) GetProjectMembers(projectID string) ([]Member, error) {
	resp, err := resty.R().
		SetAuthToken(client.token).
		Get(client.serverURL + "/v3/projects/" + projectID + "/projectroletemplatebindings")
	if err != nil {
		logrus.Errorf("Failed to query Rancher project groups: %v", err)
		return nil, err
	}

	var members []Member
	body := string(resp.Body()[:])
	memberResults := gjson.Get(body, "data")
	for _, item := range memberResults.Array() {
		userPrincipalID := item.Get("userPrincipalId").String()
		groupPrincipalID := item.Get("groupPrincipalId").String()

		newMember := Member{
			ID:             item.Get("id").String(),
			RoleTemplateID: item.Get("roleTemplateId").String(),
		}
		if userPrincipalID != "" {
			newMember.Type = MemberTypeUser
			newMember.PrincipalID = userPrincipalID
		} else {
			newMember.Type = MemberTypeGroup
			newMember.PrincipalID = groupPrincipalID

		}
		members = append(members, newMember)
	}
	return members, nil
}

func (client defaultClient) GetProjectDetail(projectID string) (*Project, error) {
	rq, err := client.GetProjectQuotas(projectID)
	if err != nil {
		return nil, err
	}

	members, err := client.GetProjectMembers(projectID)
	if err != nil {
		return nil, err
	}

	resp, err := resty.R().
		SetAuthToken(client.token).
		Get(client.serverURL + "/v3/projects/" + projectID)
	if err != nil {
		logrus.Errorf("Failed to query Rancher project groups: %v", err)
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("GetProjectDetail() failed: statusCode=%d", resp.StatusCode())
	}

	body := string(resp.Body()[:])
	return &Project{
		ID:                  projectID,
		Name:                gjson.Get(body, "name").String(),
		Description:         gjson.Get(body, "description").String(),
		ResourceQuotas:      *rq,
		Members:             members,
		PodSecurityPolicyID: gjson.Get(body, "podSecurityPolicyTemplateId").String(),
	}, nil
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
		logrus.Errorf("Failed to query Rancher namespaces: %v", err)
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
		logrus.Errorf("Failed to query Rancher projects: %v", err)
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
		logrus.Errorf("Failed to query Rancher projects: %v", err)
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
		logrus.Errorf("Failed to query Rancher project groups: %v", err)
		return nil, err
	}
	body := string(resp.Body()[:])
	groupPrincipalIds := parseValues(body, "data.#.groupPrincipalId")

	var groups []string
	for _, id := range groupPrincipalIds {
		if id != "" {
			g, err := parseGroupFromPrincipalID(id)
			if err != nil {
				logrus.Errorf("Invalid principalID '%s': %v", id, err)
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
		logrus.Errorf("Failed to query Rancher clusters: %v", err)
		return nil, err
	}
	body := string(resp.Body()[:])
	clusters := parseEntities(body, "data")
	return clusters, nil
}

func (client defaultClient) GetProjectQuotas(projectID string) (*ProjectQuotas, error) {
	resp, err := resty.R().
		SetAuthToken(client.token).
		Get(client.serverURL + "/v3/projects/" + projectID)
	if err != nil {
		logrus.Errorf("Failed to query Rancher project '%s': %v", projectID, err)
		return nil, err
	}
	body := string(resp.Body()[:])
	projectLimitsResult := gjson.Get(body, "resourceQuota.limit")
	namespaceLimitsResult := gjson.Get(body, "namespaceDefaultResourceQuota.limit")

	pq := ProjectQuotas{
		Project:   make(Quotas),
		Namespace: make(Quotas),
	}
	projectLimitsResult.ForEach(func(key, value gjson.Result) bool {
		pq.Project[key.String()] = value.String()
		return true
	})
	namespaceLimitsResult.ForEach(func(key, value gjson.Result) bool {
		pq.Namespace[key.String()] = value.String()
		return true
	})
	return &pq, nil
}

func (client defaultClient) CreateProject(clusterID string, project Project) (projectID string, err error) {
	logrus.Debugf("Creating project '%s' in cluster '%s', server-url='%s'", project.Name, clusterID, client.serverURL)
	// 	Send payload to https://rancher.example.com/v3/project?_replace=true

	logrus.Debugf("Project object: %+v", project)
	payload := map[string]interface{}{
		"type":                        "project",
		"name":                        project.Name,
		"clusterId":                   clusterID,
		"podSecurityPolicyTemplateId": project.PodSecurityPolicyID,
		"description":                 project.Description,
		"resourceQuota": map[string]interface{}{
			"limit": project.ResourceQuotas.Project,
		},
		"namespaceDefaultResourceQuota": map[string]interface{}{
			"limit": project.ResourceQuotas.Namespace,
		},
	}

	resp, err := resty.R().
		SetAuthToken(client.token).
		SetBody(payload).
		Post(client.serverURL + "/v3/project?_replace=true")
	if err != nil {
		logrus.Errorf("Failed to create project: %v", err)
		return projectID, err
	}
	body := string(resp.Body()[:])

	if resp.StatusCode() != http.StatusCreated {
		logrus.Errorf("Failed to create project: status code=%d, response=%s", resp.StatusCode(), body)
		return projectID, fmt.Errorf("created project failed, statuscode=%d", resp.StatusCode())
	}

	projectID = gjson.Get(body, "id").String()
	logrus.Debugf("Created project with ID='%s'", projectID)

	if projectID == "" {
		return "", errors.New("created projectID not found")
	}

	// Set members to project
	logrus.Debugf("Setting project members to project '%s'", projectID)
	for _, m := range project.Members {
		err := client.AddProjectMember(projectID, m)
		if err != nil {
			logrus.Errorf("Failed to bind member '%v' to project '%s': %v", m, projectID, err)
			// continue
		}
	}
	logrus.Debugf("Setting PSP '%s' to project '%s'", project.PodSecurityPolicyID, projectID)
	err = client.SetProjectPSP(projectID, project.PodSecurityPolicyID)
	return projectID, err
}

func (client defaultClient) UpdateProject(clusterID, projectID string, project Project) error {
	logrus.Debugf("Updating project ID='%s' in cluster ID='%s', server-url='%s'", projectID, clusterID, client.serverURL)
	oldPrj, err := client.GetProjectDetail(projectID)
	if err != nil {
		return err
	}

	payload := map[string]interface{}{
		"type":                        "project",
		"id":                          projectID,
		"name":                        project.Name,
		"clusterId":                   clusterID,
		"podSecurityPolicyTemplateId": project.PodSecurityPolicyID,
		"description":                 project.Description,
		"resourceQuota": map[string]interface{}{
			"limit": project.ResourceQuotas.Project,
		},
		"namespaceDefaultResourceQuota": map[string]interface{}{
			"limit": project.ResourceQuotas.Namespace,
		},
	}

	resp, err := resty.R().
		SetAuthToken(client.token).
		SetBody(payload).
		Put(client.serverURL + "/v3/projects/" + projectID + "?_replace=true")
	if err != nil {
		logrus.Errorf("Failed to update project: %v", err)
		return err
	}
	logrus.Debugf("Update project response: %v", string(resp.Body()[:]))

	if oldPrj.PodSecurityPolicyID != project.PodSecurityPolicyID {
		// update PSP
		logrus.Debugf("Project PSP changed, updating to '%s'", project.PodSecurityPolicyID)
		err = client.SetProjectPSP(projectID, project.PodSecurityPolicyID)
		if err != nil {
			return err
		}
	}

	var deletedMembers, newMembers []Member
	for _, m := range project.Members {
		if !hasMember(oldPrj.Members, m) {
			newMembers = append(newMembers, m)
		}
	}
	for _, m := range oldPrj.Members {
		if !hasMember(project.Members, m) {
			deletedMembers = append(deletedMembers, m)
		}
	}
	logrus.Debugf("New members: %v", newMembers)
	logrus.Debugf("Deleted members: %v", deletedMembers)

	for _, m := range newMembers {
		err = client.AddProjectMember(projectID, m)
		if err != nil {
			logrus.Errorf("Adding member failed: %v", err)
		}
	}

	for _, m := range deletedMembers {
		err = client.DeleteProjectMember(m.ID)
		if err != nil {
			logrus.Errorf("Deleting member failed: %v", err)
		}
	}

	logrus.Debugf("Updated project with ID='%s'", projectID)

	return nil
}

// hasMember returns true if the target is in the array
func hasMember(members []Member, target Member) bool {
	for _, m := range members {
		if m.Compare(target) {
			return true
		}
	}
	return false
}

func (client defaultClient) AddProjectMember(projectID string, member Member) (err error) {
	payload := map[string]interface{}{
		"type":                  "projectRoleTemplateBinding",
		"subjectKind":           member.Type,
		"userId":                "",
		"projectRoleTemplateId": "",
		"projectId":             projectID,
		"groupPrincipalId":      "",
		"userPrincipalId":       "",
		"roleTemplateId":        member.RoleTemplateID,
	}

	switch member.Type {
	case MemberTypeUser:
		payload["userPrincipalId"] = member.PrincipalID
		break
	case MemberTypeGroup:
		payload["groupPrincipalId"] = member.PrincipalID
		break
	default:
		return errors.New("invalid member type")
	}

	resp, err := resty.R().
		SetAuthToken(client.token).
		SetBody(payload).
		Post(client.serverURL + "/v3/projectroletemplatebinding")
	if err != nil {
		return err
	}

	logrus.Debugf("Binding role response: %v", string(resp.Body()[:]))

	if resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("add project member failed, code=%d", resp.StatusCode())
	}
	return nil
}

func (client defaultClient) SetProjectPSP(projectID string, PodSecurityPolicyID string) error {
	payload := map[string]interface{}{
		"podSecurityPolicyTemplateId": PodSecurityPolicyID,
	}
	resp, err := resty.R().
		SetAuthToken(client.token).
		SetBody(payload).
		Post(client.serverURL + "/v3/projects/" + projectID + "?action=setpodsecuritypolicytemplate")
	if err != nil {
		return err
	}
	logrus.Debugf("Set ProjectPSP response: %v", string(resp.Body()[:]))
	return nil
}

func (client defaultClient) DeleteProjectMember(ID string) error {
	logrus.Debugf("Deleting member %s", ID)
	resp, err := resty.R().
		SetAuthToken(client.token).
		Delete(client.serverURL + "/v3/projectRoleTemplateBindings/" + ID)
	if err != nil {
		return err
	}
	statusCode := resp.StatusCode()
	if statusCode != http.StatusOK {
		logrus.Errorf("Delete project member failed, response: %v", string(resp.Body()[:]))
		return fmt.Errorf("delete failed, response=%d", statusCode)
	}
	return nil
}
