package ansible

import (
	"bytes"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/samber/lo"
	"github.com/sosedoff/ansible-vault-go"
	"gopkg.in/yaml.v3"
)

var (
	vaultHeader  = []byte("$ANSIBLE_VAULT;1.1;AES256")
	variableDirs = []string{"group_vars", "host_vars", "defaults", "vars"}
)

type Result struct {
	Path     string
	Variable string
	Value    yaml.Node
}

func decryptVariable(v *yaml.Node, password string) {
	if bytes.HasPrefix([]byte(v.Value), vaultHeader) {
		decrypted, err := vault.Decrypt(v.Value, password)
		if err != nil {
			log.Fatal(err)
		}
		v.Value = decrypted
	}
}

func Find(root, password, variable string) ([]Result, error) {
	var results []Result

	err := walk(root, password, func(path string, yml map[string]yaml.Node) {
		if v, found := yml[variable]; found {
			decryptVariable(&v, password)
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

	err = walk(root, password, func(path string, yml map[string]yaml.Node) {
		for k, v := range yml {
			if rgx.MatchString(k) {
				decryptVariable(&v, password)
				results = append(results, Result{Path: path, Variable: k, Value: v})
			}
		}
	})
	return results, err
}

func walk(root, password string, run func(path string, yml map[string]yaml.Node)) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && d.Name() != root && strings.HasPrefix(d.Name(), ".") {
			return filepath.SkipDir
		}

		parts := strings.Split(path, "/")
		if !lo.ContainsBy(variableDirs, func(dir string) bool { return lo.Contains(parts, dir) }) {
			return nil
		}
		if !lo.Contains([]string{".yaml", ".yml"}, filepath.Ext(d.Name())) {
			return nil
		}

		yml, err := parseFile(path, password)
		if err != nil {
			return err
		}
		run(path, yml)
		return nil
	})
}

func parseFile(path, password string) (map[string]yaml.Node, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if bytes.HasPrefix(content, vaultHeader) {
		decrypted, err := vault.Decrypt(string(content), password)
		if err != nil {
			return nil, err
		}
		content = []byte(decrypted)
	}

	yml := make(map[string]yaml.Node)
	err = yaml.Unmarshal(content, &yml)
	return yml, err
}
