package ansible

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/samber/lo"
	"github.com/sosedoff/ansible-vault-go"
	"gopkg.in/yaml.v2"
)

var (
	vaultHeader  = []byte("$ANSIBLE_VAULT;1.1;AES256")
	variableDirs = []string{"/group_vars/", "/host_vars/", "/defaults/", "/vars/"}
)

type Result struct {
	Path     string
	Variable string
	Value    any
}

func Find(root, password, variable string) ([]Result, error) {
	var results []Result

	err := walk(root, password, func(path string, yml map[any]any, isVault bool) {
		if v, found := yml[variable]; found {
			results = append(results, Result{Path: path, Variable: variable, Value: v})
		}
	})
	return results, err
}

func FindRegex(root, password, variable string) ([]Result, error) {
	var results []Result
	rgx, err := regexp.Compile(variable)
	if err != nil {
		return nil, err
	}

	err = walk(root, password, func(path string, yml map[any]any, isVault bool) {
		for k, v := range yml {
			if rgx.MatchString(k.(string)) {
				results = append(results, Result{Path: path, Variable: k.(string), Value: v})
			}
		}
	})
	return results, err
}

func walk(root, password string, run func(path string, yml map[any]any, isVault bool)) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && d.Name() != root && strings.HasPrefix(d.Name(), ".") {
			return filepath.SkipDir
		}
		if !lo.ContainsBy(variableDirs, func(dir string) bool { return strings.Contains(path, dir) }) {
			return nil
		}
		if !lo.Contains([]string{".yaml", ".yml"}, filepath.Ext(d.Name())) {
			return nil
		}

		yml, isVault, err := parseFile(path, password)
		if err != nil {
			return err
		}
		run(path, yml, isVault)
		return nil
	})
}

func parseFile(path, password string) (map[any]any, bool, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, false, err
	}

	isVault := bytes.HasPrefix(content, vaultHeader)
	if isVault {
		decrypted, err := vault.Decrypt(string(content), password)
		if err != nil {
			return nil, isVault, err
		}
		content = []byte(decrypted)
	}

	yml := make(map[any]any)
	err = yaml.Unmarshal(content, &yml)
	return yml, isVault, err
}
