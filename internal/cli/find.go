package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/jonasbleyl/ansible-vars/pkg/ansible"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	startBlueOutput = "\033[34m"
	stopColorOutput = "\033[0m"
)

var (
	showFileNamesOnly bool
	regex             bool
)

func setupFind(cmd *cobra.Command) {
	findCmd := &cobra.Command{
		Use:   "find VARIABLE [DIRECTORY]",
		Short: "Find where ansible variables are defined",
		Args:  cobra.RangeArgs(1, 2),
		RunE:  runFind,
	}

	findCmd.Flags().BoolVarP(&showFileNamesOnly, "files-with-matches", "l", false, "Only print the filenames of matching files")
	findCmd.Flags().BoolVarP(&regex, "regex", "r", false, "Use regex to search for variables")

	cmd.AddCommand(findCmd)
}

func runFind(cmd *cobra.Command, args []string) error {
	dir := "."
	if len(args) > 1 {
		dir = args[1]
	}

	contents, err := os.ReadFile(vaultFile)
	if err != nil {
		return err
	}
	password := strings.TrimSuffix(string(contents), "\n")

	results, err := executeFind(args[0], dir, password)
	if err != nil {
		return err
	}
	outputFind(cmd, results)
	return nil
}

func executeFind(variable, dir, password string) ([]ansible.Result, error) {
	if regex {
		return ansible.FindRegex(dir, password, variable)
	}
	return ansible.Find(dir, password, variable)
}

func outputFind(cmd *cobra.Command, results []ansible.Result) {
	for _, r := range results {
		cmd.Println(fmt.Sprintf("%s%s%s", startBlueOutput, r.Path, stopColorOutput))

		if !showFileNamesOnly {
			output := make(map[string]yaml.Node)
			output[r.Variable] = r.Value

			yml, _ := yaml.Marshal(&output)
			cmd.Print(string(yml))
		}
	}
}
