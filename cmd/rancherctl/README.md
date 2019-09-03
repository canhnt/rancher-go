# Manage projects with Rancher command line
This project creates `rancherctl` command, that can manage projects in a K8s cluster controlled by Rancher application server.

Version: 0.1.0

## Usage
### Help
```
$ rancherctl
Rancher Project CLI, managing Rancher projects
Usage: rancherctl [OPTIONS] COMMAND [arg...]
Version: dev

Options:
  --debug              Debug logging
  --rancher-url value  URL of the Rancher server
  --cluster value      Target cluster ID to manage projects
  --token value        Security token used to access Rancher APIs
  --help, -h           show help
  --version, -v        print the version

Commands:
  ls         List projects
  get        Get project
  apply      Create or update multiple projects
  delete     Remove a project
  help, [h]  Shows a list of commands or help for one command

Run 'rancherctl COMMAND --help' for more information on a command.
```

### List projects
```
TOKEN=`token-abcd:123456...'
CLUSTER_ID='c-a1bcd'
$ rancherctl --rancher-url=https://rancher.example.org --cluster=${CLUSTER_ID} --token=${TOKEN} ls
Projects in cluster 'c-a1bcd'
ID 			 Name
c-a1bcd:p-242lq 	 Customer_Master
c-a1bcd:p-27z28 	 WebServer
c-a1bcd:p-28vtz 	 MyProject1
c-a1bcd:p-29bs7 	 Sandbox
c-a1bcd:p-29qn8 	 BackendApp
c-a1bcd:p-2kks9 	 Workflows
c-a1bcd:p-2l2l4 	 Playground
c-a1bcd:p-2wfqv    canhnt
```

### Get project detail
```
$ rancherctl --rancher-url=https://rancher.example.org --cluster=${CLUSTER_ID} --token=${TOKEN} get c-a1bcd:p-2wfqv
id: c-a1bcd:p-2wfqv
name: canhnt
description: Personal project of Canh
podSecurityPolicyId: mypsp
members:
- id: p-mrjs8:prtb-b42z4
  type: User
  principalId: openldap_user://cn=canhnt,ou=People,dc=example
  roleTemplateId: project-owner
projectQuotas:
  project:
    limitsCpu: 2000m
    limitsMemory: 4096Mi
    requestsCpu: 1000m
    requestsMemory: 2048Mi
    requestsStorage: 80Gi
    type: /v3/schemas/resourceQuotaLimit
  namespace:
    limitsCpu: 500m
    limitsMemory: 1024Mi
    requestsCpu: 200m
    requestsMemory: 512Mi
    requestsStorage: 20Gi
    type: /v3/schemas/resourceQuotaLimit
```

### Create projects
- Create a YAML file `projects.yaml` containing projects configuration
```yaml
projects:
  - id: 'c-a1b2c3:p-b1c2d3'
    name: "demo-project1"
    description: "1st project storing in YAML file"
    podSecurityPolicyId: mypsp
    members:
      - type: Group
        principalId: openldap_group://cn=developers,ou=Groups,dc=example
        roleTemplateId: project-member
      - type: Group
        principalId: openldap_group://cn=devops,ou=Groups,dc=example
        roleTemplateId: project-member
      - type: User
        principalId: openldap_user://cn=canhnt,ou=People,dc=example
        roleTemplateId: project-owner
      - type: Group
        principalId: "openldap_group://cn=all_developers_all,ou=Groups,dc=example"
        roleTemplateId: read-only       
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
    podSecurityPolicyId: mypsp
    members:
      - type: Group
        principalId: openldap_group://cn=developers,ou=Groups,dc=example
        roleTemplateId: project-member
      - type: Group
        principalId: openldap_group://cn=devops,ou=Groups,dc=example
        roleTemplateId: project-member
      - type: User
        principalId: openldap_user://cn=canhnt,ou=People,dc=example
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
```

- Create defined projects:
```
rancherctl --rancher-url=https://rancher.example.org --cluster=${CLUSTER_ID} --token=${TOKEN} apply --filename projects.yaml
```
