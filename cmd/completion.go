package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion script",
	Long: `Generate shell completion script for crestic.

To load completions:

Bash:

  $ source <(crestic completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ crestic completion bash > /etc/bash_completion.d/crestic
  # macOS:
  $ crestic completion bash > $(brew --prefix)/etc/bash_completion.d/crestic

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ crestic completion zsh > "${fpath[1]}/_crestic"

  # You will need to start a new shell for this setup to take effect.

Fish:

  $ crestic completion fish | source

  # To load completions for each session, execute once:
  $ crestic completion fish > ~/.config/fish/completions/crestic.fish

PowerShell:

  PS> crestic completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> crestic completion powershell > crestic.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			return cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			return cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
