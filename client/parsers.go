package client

import (
	"errors"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/tidwall/gjson"
)

// Parse 'openldap_group://cn=foo,ou=Groups,dc=example.com' to 'foo'
func parseGroupFromPrincipalID(principalID string) (string, error) {
	if !strings.HasPrefix(principalID, "openldap_group") {
		return "", errors.New("Invalid LDAP principalID")
	}

	dn := principalID[len("openldap_group://"):]
	dnParts := strings.Split(dn, ",")
	if len(dnParts) < 1 {
		return "", errors.New("CN in the principalID not found")
	}
	cn := dnParts[0]
	cnParts := strings.Split(cn, "=")
	if len(cnParts) < 2 {
		return "", errors.New("Invalid CN value in the principalID")
	}
	if cnParts[0] != "cn" {
		return "", errors.New("CN field not found")
	}
	return cnParts[1], nil
}

func parseValues(jsonData string, jsonPath string) []string {
	var values []string
	results := gjson.Get(jsonData, jsonPath).Array()
	for _, val := range results {
		values = append(values, val.String())
	}
	return values
}

// parseEntities extracts 'id' and 'name' attributes of the json array objects
func parseEntities(jsonData string, jsonPath string) []Entity {
	var entities []Entity
	result := gjson.Get(jsonData, jsonPath)
	result.ForEach(func(key, value gjson.Result) bool {
		name := value.Get("name").String()
		id := value.Get("id").String()
		if name == "" || id == "" {
			logrus.Errorf("Either cluster name or id is empty: name='%s', id='%s'", name, id)
			return true // continue next item
		}
		entities = append(entities, Entity{ID: id, Name: name})
		return true
	})
	return entities
}
