package ansible

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
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
		{Path: testDataDir + "/group_vars/values.yaml", Variable: "test_var", Value: yaml.Node{Value: "this value is encrypted"}},
	}
	results, err := Find(testDataDir, vaultPassword, "test_var")
	assert.NoError(t, err)
	assertResultsEqual(t, want, results)
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
		{Path: testDataDir + "/defaults/vault.yaml", Variable: "test_var2", Value: yaml.Node{Value: "value"}},
		{Path: testDataDir + "/group_vars/vault.yaml", Variable: "test_var", Value: yaml.Node{Value: "value"}},
		{Path: testDataDir + "/host_vars/vault.yaml", Variable: "test_var", Value: yaml.Node{Value: "value"}},
		{Path: testDataDir + "/inventories/host_vars/vault.yaml", Variable: "test_var", Value: yaml.Node{Value: "value"}},
		{Path: testDataDir + "/vars/vars.yml", Variable: "test_var", Value: yaml.Node{Value: "value"}},
		{Path: testDataDir + "/group_vars/values.yaml", Variable: "test_var", Value: yaml.Node{Value: "this value is encrypted"}},
	}
	results, err := FindRegex(testDataDir, vaultPassword, "test_.*")
	assert.NoError(t, err)
	assertResultsEqual(t, want, results)
}

func TestFindRegex_badRegex(t *testing.T) {
	_, err := FindRegex(testDataDir, vaultPassword, "*")
	assert.Error(t, err)
}

func assertResultsEqual(t *testing.T, want, got []Result) {
	t.Helper()

	assert.Len(t, got, len(want))

	sort.Slice(want, func(i, j int) bool {
		if want[i].Path != want[j].Path {
			return want[i].Path < want[j].Path
		}
		return want[i].Variable < want[j].Variable
	})

	sort.Slice(got, func(i, j int) bool {
		if got[i].Path != got[j].Path {
			return got[i].Path < got[j].Path
		}
		return got[i].Variable < got[j].Variable
	})

	for i, r := range got {
		assert.Equal(t, want[i].Path, r.Path)
		assert.Equal(t, want[i].Variable, r.Variable)
		assert.Equal(t, want[i].Value.Value, r.Value.Value)
	}
}
