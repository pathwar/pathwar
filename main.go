package main // import "pathwar.pw"

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"pathwar.pw/server"
)

func main() {
	rootCmd := newRootCommand()
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "pathwar.pw",
	}
	cmd.PersistentFlags().BoolP("help", "h", false, "print usage")
	//cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose mode")

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// setup logging etc
		return nil
	}

	cmd.AddCommand(
		server.NewServerCommand(),
	)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	return cmd
}
