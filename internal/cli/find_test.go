package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	tests := map[string]struct {
		args string
		want string
	}{
		"default": {
			args: "test_var ../../testdata/defaults --vault ../../testdata/.vault",
			want: startBlueOutput + "../../testdata/defaults/vault.yaml" + stopColorOutput + "\ntest_var: value\n",
		},
		"filenames only": {
			args: "test_var ../../testdata/defaults --vault ../../testdata/.vault --files-with-matches",
			want: "../../testdata/defaults/vault.yaml\n",
		},
		"regex": {
			args: "--regex test_.* ../../testdata/defaults --vault ../../testdata/.vault",
			want: startBlueOutput + "../../testdata/defaults/vault.yaml" + stopColorOutput + "\ntest_var: value\n",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			out, err := testExecute(tc.args)
			assert.NoError(t, err)
			assert.Equal(t, tc.want, out)
		})
	}
}

func TestFind_error(t *testing.T) {
	tests := map[string]struct {
		args string
	}{
		"invalid directory":  {args: "test_var /does/not/exist --vault ../../testdata/.vault"},
		"invalid vault file": {args: "test_var ../../testdata/group_vars --vault does/not/exist"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := testExecute(tc.args)
			assert.Error(t, err)
		})
	}
}

func testExecute(args string) (string, error) {
	out := new(bytes.Buffer)
	cmd := Setup()
	cmd.SetOut(out)
	cmd.SetErr(out)
	cmd.SetArgs(strings.Split(args, " "))
	err := cmd.Execute()
	return out.String(), err
}
