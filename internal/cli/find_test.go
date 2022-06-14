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
			args: "find test_var ../../testdata/defaults --vault ../../testdata/.vault",
			want: startBlueOutput + "../../testdata/defaults/vault.yaml" + stopColorOutput + "\ntest_var: value\n",
		},
		"filenames only": {
			args: "find test_var ../../testdata/defaults --vault ../../testdata/.vault --files-with-matches",
			want: startBlueOutput + "../../testdata/defaults/vault.yaml" + stopColorOutput + "\n",
		},
		"regex": {
			args: "find --regex test_.* ../../testdata/defaults --vault ../../testdata/.vault",
			want: startBlueOutput + "../../testdata/defaults/vault.yaml" + stopColorOutput + "\ntest_var: value\n",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			out, err := execute(tc.args)
			assert.NoError(t, err)
			assert.Equal(t, tc.want, out)
		})
	}
}

func TestFind_error(t *testing.T) {
	tests := map[string]struct {
		args string
	}{
		"invalid directory":  {args: "find var /does/not/exist --vault ../../testdata/.vault"},
		"invalid vault file": {args: "find var ../../testdata/group_vars --vault does/not/exist"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := execute(tc.args)
			assert.Error(t, err)
		})
	}
}

func execute(args string) (string, error) {
	out := new(bytes.Buffer)
	rootCmd := Setup()
	rootCmd.SetOut(out)
	rootCmd.SetErr(out)
	rootCmd.SetArgs(strings.Split(args, " "))
	err := rootCmd.Execute()
	return out.String(), err
}
