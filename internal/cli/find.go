package cli

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/jonasbleyl/ansible-find/pkg/ansible"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	startBlueOutput = "\033[34m"
	stopColorOutput = "\033[0m"
)

var (
	vaultFile         string
	showFileNamesOnly bool
	regex             bool
)

func Setup() *cobra.Command {
	cmd := &cobra.Command{
		Use: "ansible-find VARIABLE [DIRECTORY]",
		Long: `A CLI tool to find where ansible variables are defined.

This tool will only search YAML files that reside within the
following directories: [group_vars, host_vars, defaults, vars]`,
		Args: cobra.RangeArgs(1, 2),
		RunE: run,
	}

	cmd.Flags().StringVarP(&vaultFile, "vault", "v", ".vault", "ansible vault password file")
	cmd.Flags().BoolVarP(&showFileNamesOnly, "files-with-matches", "l", false, "only print the filenames of matching files")
	cmd.Flags().BoolVarP(&regex, "regex", "r", false, "use regex to search for variables")

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	dir := "."
	if len(args) > 1 {
		dir = args[1]
	}

	contents, err := os.ReadFile(vaultFile)
	if err != nil {
		return err
	}
	password := strings.TrimSuffix(string(contents), "\n")

	results, err := execute(args[0], dir, password)
	if err != nil {
		return err
	}
	return output(cmd, results)
}

func execute(variable, dir, password string) ([]ansible.Result, error) {
	if regex {
		return ansible.FindRegex(dir, password, variable)
	}
	return ansible.Find(dir, password, variable)
}

func output(cmd *cobra.Command, results []ansible.Result) error {
	for _, r := range results {
		if showFileNamesOnly {
			cmd.Println(r.Path)
			continue
		}

		var b bytes.Buffer
		encoder := yaml.NewEncoder(&b)
		encoder.SetIndent(2)

		yml := map[string]yaml.Node{r.Variable: r.Value}
		err := encoder.Encode(&yml)
		if err != nil {
			return err
		}

		cmd.Println(fmt.Sprintf("%s%s%s", startBlueOutput, r.Path, stopColorOutput))
		cmd.Print(b.String())
	}
	return nil
}
