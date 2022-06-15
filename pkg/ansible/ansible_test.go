package ansible

import (
	"gopkg.in/yaml.v3"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testDataDir   = "../../testdata"
	vaultPassword = "pa$$word"
)

func TestFind(t *testing.T) {
	want := []Result{
		{Path: testDataDir + "/defaults/vault.yaml", Variable: "test_var", Value: yaml.Node{Value: "value"}},
		{Path: testDataDir + "/group_vars/vault.yaml", Variable: "test_var", Value: yaml.Node{Value: "value"}},
		{Path: testDataDir + "/host_vars/vault.yaml", Variable: "test_var", Value: yaml.Node{Value: "value"}},
		{Path: testDataDir + "/inventories/host_vars/vault.yaml", Variable: "test_var", Value: yaml.Node{Value: "value"}},
		{Path: testDataDir + "/vars/vars.yml", Variable: "test_var", Value: yaml.Node{Value: "value"}},
	}
	results, err := Find(testDataDir, vaultPassword, "test_var")
	assert.NoError(t, err)
	assert.Len(t, results, len(want))

	for i, r := range results {
		assert.Equal(t, want[i].Path, r.Path)
		assert.Equal(t, want[i].Variable, r.Variable)
		assert.Equal(t, want[i].Value.Value, r.Value.Value)
	}
}

func TestFind_error(t *testing.T) {
	tests := map[string]struct {
		path     string
		password string
	}{
		"non-existing path":    {path: "does/not/exist", password: vaultPassword},
		"wrong vault password": {path: testDataDir + "/defaults/vault.yaml", password: "wrong password"},
		"non-existing file":    {path: testDataDir + "/defaults/does_not_exist.yaml", password: vaultPassword},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := Find(tc.path, tc.password, "test_var")
			assert.Error(t, err)
		})
	}
}

func TestFindRegex(t *testing.T) {
	want := []Result{
		{Path: testDataDir + "/defaults/vault.yaml", Variable: "test_var", Value: yaml.Node{Value: "value"}},
		{Path: testDataDir + "/group_vars/vault.yaml", Variable: "test_var", Value: yaml.Node{Value: "value"}},
		{Path: testDataDir + "/host_vars/vault.yaml", Variable: "test_var", Value: yaml.Node{Value: "value"}},
		{Path: testDataDir + "/inventories/host_vars/vault.yaml", Variable: "test_var", Value: yaml.Node{Value: "value"}},
		{Path: testDataDir + "/vars/vars.yml", Variable: "test_var", Value: yaml.Node{Value: "value"}},
	}
	results, err := FindRegex(testDataDir, vaultPassword, "test_.*")
	assert.NoError(t, err)

	for i, r := range results {
		assert.Equal(t, want[i].Path, r.Path)
		assert.Equal(t, want[i].Variable, r.Variable)
		assert.Equal(t, want[i].Value.Value, r.Value.Value)
	}
}

func TestFindRegex_badRegex(t *testing.T) {
	_, err := FindRegex(testDataDir, vaultPassword, "*")
	assert.Error(t, err)
}
