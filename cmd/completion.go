package cmd

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

// generalComDesc describes the long description for the completion command
const generalComDesc = `
Generate the autocompletion script for medusa for the specific shell.

Bash:
To load completions in the current shell session:

    source <(medusa completion bash)

To load completions for every new session, execute once:
- Linux:
      medusa completion bash > /etc/bash_completion.d/medusa

- macOS:
      medusa completion bash > /usr/local/etc/bash_completion.d/medusa

Zsh:
To load completions in the current shell session:

    source <(medusa completion zsh)

To load completions for every new session, execute once:

    medusa completion zsh > "${fpath[1]}/_medusa"

PowerShell:
To load completions in the current shell session:
PS> medusa completion powershell | Out-String | Invoke-Expression

To load completions for every new session, run:
PS> medusa completion powershell > medusa.ps1
and source this file from your PowerShell profile.
`

var supportedShells = []string{"bash", "zsh", "powershell"}

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:           "completion <shell>",
	Short:         "Generate the autocompletion script for medusa for the specific shell",
	Long:          generalComDesc,
	Args:          cmdValidateCompletionArgs,
	RunE:          cmdRunCompletion,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.AddCommand(completionCmd)
}

// cmdValidateCompletionArgs validates CLI arguments
func cmdValidateCompletionArgs(cmd *cobra.Command, args []string) error {
	// Make sure we have exactly 1 argument
	if err := cobra.ExactArgs(1)(cmd, args); err != nil {
		err = fmt.Errorf("completion requires only 1 shell argument (options: %s)", strings.Join(supportedShells, ", "))
		cmdLogger.Error("Failed to validate args for completion command", err)
		return err
	}

	// Make sure that the shell is a supported type
	if contains := slices.Contains(supportedShells, args[0]); !contains {
		err := fmt.Errorf("%s is not a supported shell", args[0])
		cmdLogger.Error("Failed to validate args for completion command", err)
		return err
	}

	return nil
}

// cmdRunCompletion executes the completion CLI command
func cmdRunCompletion(cmd *cobra.Command, args []string) error {
	// NOTE: Please be aware that if the supported shells changes, then this switch statement must also change
	var err error
	switch args[0] {
	case "bash":
		err = cmd.Root().GenBashCompletionV2(os.Stdout, true)
	case "zsh":
		err = cmd.Root().GenZshCompletion(os.Stdout)
	case "powershell":
		err = cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
	default:
		// We are throwing a panic here because our validation function should have handled this and something is wrong.
		cmdLogger.Panic("Failed to run the completion command", fmt.Errorf("%s is not a supported shell type", args[0]))
	}

	// Log an error if we encountered one
	if err != nil {
		cmdLogger.Error("Failed to run the completion command", err)
	}
	return err
}
