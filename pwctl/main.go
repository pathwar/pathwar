package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	rootCmd := &cobra.Command{
		Use: os.Args[0],
	}
	rootCmd.PersistentFlags().BoolP("help", "h", false, "print usage")
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		return nil
	}
	rootCmd.AddCommand(&cobra.Command{
		Use:  "entrypoint",
		Args: cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			// FIXME: lock to block other commands

			// prepare the level
			cmd := exec.Command("/pathwar-hooks/on-init")
			if err := cmd.Run(); err != nil {
				return err
			}

			// FIXME: add a self-destruct mode that allow having root access only at runtime

			// switch to original's entrypoint
			binary, err := exec.LookPath(args[0])
			if err != nil {
				return err
			}
			env := os.Environ()
			if err := syscall.Exec(binary, args, env); err != nil {
				return err
			}
			return nil
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use: "env",
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, line := range os.Environ() {
				fmt.Println(line)
			}
			return nil
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use: "config",
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := getConfig()
			if err != nil {
				return err
			}
			out, _ := json.Marshal(config)
			fmt.Println(string(out))
			return nil
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use: "passphrase",
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := getConfig()
			if err != nil {
				return err
			}
			// FIXME: parse index from CLI
			fmt.Println(config.Passphrases[0])
			return nil
		},
	})
	args := os.Args[1:]
	if args[0] == "entrypoint" && len(args) > 1 && args[1] != "--" {
		args = append([]string{"entrypoint", "--"}, args[1:]...)
	}
	rootCmd.SetArgs(args)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
