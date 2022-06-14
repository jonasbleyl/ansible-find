package ansible

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testDataDir   = "../../testdata"
	vaultPassword = "pa$$word"
)

func TestFind(t *testing.T) {
	tests := map[string]struct {
		path     string
		password string
		variable string
		want     []Result
	}{
		"find unencrypted variable": {
			path:     testDataDir,
			password: vaultPassword,
			variable: "unencrypted_var",
			want: []Result{
				{Path: testDataDir + "/vars/vars.yml", Variable: "unencrypted_var", Value: "value"},
			},
		},
		"find encrypted variable": {
			path:     testDataDir,
			password: vaultPassword,
			variable: "encrypted_var",
			want: []Result{
				{Path: testDataDir + "/defaults/vault.yaml", Variable: "encrypted_var", Value: "value"},
				{Path: testDataDir + "/group_vars/vault.yaml", Variable: "encrypted_var", Value: "value"},
				{Path: testDataDir + "/host_vars/vault.yaml", Variable: "encrypted_var", Value: "value"},
				{Path: testDataDir + "/inventories/host_vars/vault.yaml", Variable: "encrypted_var", Value: "value"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			output, err := Find(tc.path, tc.password, tc.variable)
			assert.NoError(t, err)
			assert.ElementsMatch(t, tc.want, output)
		})
	}
}

func TestFindError(t *testing.T) {
	tests := map[string]struct {
		path     string
		password string
	}{
		"non-existing path":    {path: "does/not/exist", password: vaultPassword},
		"wrong vault password": {path: testDataDir + "/defaults/vault.yaml", password: "wrong password"},
		"non-existing file":    {path: testDataDir + "/defaults/does_not_exist.yaml", password: vaultPassword},
		//"":      {path: testDataDir + ".skip", password: vaultPassword},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := Find(tc.path, tc.password, "encrypted_var")
			assert.Error(t, err)
		})
	}
}
