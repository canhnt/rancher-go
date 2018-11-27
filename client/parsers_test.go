package client

import (
	"testing"
)

func Test_parseGroupFromPrincipalID(t *testing.T) {
	type args struct {
		principalID string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"valid principal id",
			args{"openldap_group://cn=foo,ou=Groups,dc=example.com"},
			"foo",
			false,
		},
		{
			"invalid principalID 1",
			args{"openldap_group://ou=Groups,dc=example.com"},
			"",
			true,
		},
		{
			"invalid principalID 2",
			args{"cn=foo,ou=Groups,dc=example.com"},
			"",
			true,
		},
		{
			"invalid principalID 3",
			args{"openldap://cn=foo,ou=Groups,dc=example.com"},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseGroupFromPrincipalID(tt.args.principalID)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseGroupFromPrincipalID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseGroupFromPrincipalID() = %v, want %v", got, tt.want)
			}
		})
	}
}
