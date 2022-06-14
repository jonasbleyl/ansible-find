package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/jonasbleyl/ansible-vars/pkg/ansible"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const (
	startBlueOutput = "\033[34m"
	stopColorOutput = "\033[0m"
)

var (
	showFileNamesOnly bool
)

func setupFind(cmd *cobra.Command) {
	findCmd := &cobra.Command{
		Use:   "find VARIABLE [DIRECTORY]",
		Short: "Find where ansible variables are defined",
		Args:  cobra.RangeArgs(1, 2),
		RunE:  findRun,
	}

	findCmd.Flags().BoolVarP(&showFileNamesOnly, "files-with-matches", "l", false, "Only print the filenames of matching files")

	cmd.AddCommand(findCmd)
}

func findRun(cmd *cobra.Command, args []string) error {
	variable := args[0]
	dir := "."

	if len(args) > 1 {
		dir = args[1]
	}

	contents, err := os.ReadFile(vaultFile)
	password := strings.TrimSuffix(string(contents), "\n")
	if err != nil {
		return err
	}

	results, err := ansible.Find(dir, password, variable)
	if err != nil {
		return err
	}

	for _, r := range results {
		cmd.Println(fmt.Sprintf("%s%s%s", startBlueOutput, r.Path, stopColorOutput))

		if !showFileNamesOnly {
			output := make(map[any]any)
			output[variable] = r.Value

			yml, _ := yaml.Marshal(&output)
			cmd.Print(string(yml))
		}
	}
	return nil
}
