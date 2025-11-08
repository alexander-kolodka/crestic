package cmd

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

//go:embed config.example.yaml
var configExample string

var configExampleCmd = &cobra.Command{
	Use:   "config.example",
	Short: "Print example configuration file",
	Long:  "Outputs a comprehensive example configuration file with detailed comments and documentation",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Fprint(os.Stdout, configExample)
	},
}

func init() {
	rootCmd.AddCommand(configExampleCmd)
}
