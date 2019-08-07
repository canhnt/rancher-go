package client

type Quotas map[string]string

type ProjectQuotas struct {
	Project   Quotas `yaml:"project,omitempty"`
	Namespace Quotas `yaml:"namespace,omitempty"`
}

const (
	MemberTypeUser  = "User"
	MemberTypeGroup = "Group"
)

type Member struct {
	ID             string `yaml:"id"`
	Type           string `yaml:"type"`
	PrincipalID    string `yaml:"principalId,omitempty"`
	RoleTemplateID string `yaml:"roleTemplateId,omitempty"`
}

type Project struct {
	ID                  string        `yaml:"id"`
	Name                string        `yaml:"name"`
	Description         string        `yaml:"description,omitempty"`
	PodSecurityPolicyID string        `yaml:"podSecurityPolicyId,omitempty"`
	Members             []Member      `yaml:"members,omitempty"`
	ResourceQuotas      ProjectQuotas `yaml:"projectQuotas,omitempty"`
}

type ProjectList struct {
	Projects []Project `yaml:"projects"`
}

// compare does the comparision but ignores the ID field
func (m Member) Compare(p Member) bool {
	// TODO use https://github.com/google/go-cmp?
	return m.Type == p.Type && m.PrincipalID == p.PrincipalID && m.RoleTemplateID == p.RoleTemplateID
}
