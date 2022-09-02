package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"

	"github.com/peterbourgon/ff/v3/ffcli"
	"moul.io/banner"
	"moul.io/u"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwinit"
)

func main() {
	log.SetFlags(0)

	entrypoint := &ffcli.Command{
		Name:       "entrypoint",
		ShortUsage: "pwinit entrypoint [args...]",
		Exec: func(ctx context.Context, args []string) error {
			// FIXME: lock to block other commands

			if len(args) < 1 {
				return flag.ErrHelp
			}

			fmt.Println(banner.Inline("pwinit"))

			if !u.FileExists("/pwinit/config.json") {
				log.Print("no such config file, skipping on-init hook")
			} else {
				if u.FileExists("/pwinit/on-init") {
					log.Print("starting on-init hook")
					// prepare the challenge
					err := os.Chmod("/pwinit/on-init", 0o555)
					if err != nil {
						return errcode.ErrExecuteOnInitHook.Wrap(err)
					}
					cmd := exec.Command("/pwinit/on-init")
					err = cmd.Run()
					if err != nil {
						return errcode.ErrExecuteOnInitHook.Wrap(err)
					}
				}

				// clean pwinit config file that contains passphrases
				for _, filename := range []string{"/pwinit/config.json", "/pwinit/on-init"} {
					if !u.FileExists(filename) {
						continue
					}
					err := os.Remove(filename)
					if err != nil {
						return errcode.ErrRemoveInitConfig.Wrap(err)
					}
				}
			}

			// FIXME: add a self-destruct mode that allows having root access only at runtime

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
	}

	env := &ffcli.Command{
		Name:       "env",
		ShortUsage: "pwinit entrypoint [args...]",
		Exec: func(ctx context.Context, args []string) error {
			for _, line := range os.Environ() {
				fmt.Println(line)
			}
			return nil
		},
	}

	config := &ffcli.Command{
		Name:       "config",
		ShortUsage: "pwinit config [args...]",
		Exec: func(ctx context.Context, args []string) error {
			config, err := getConfig()
			if err != nil {
				return err
			}
			out, _ := json.Marshal(config)
			fmt.Println(string(out))
			return nil
		},
	}

	passphrase := &ffcli.Command{
		Name:       "passphrase",
		ShortUsage: "pwinit passphrase ID",
		Exec: func(ctx context.Context, args []string) error {
			if len(args) != 1 {
				return flag.ErrHelp
			}
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			config, err := getConfig()
			if err != nil {
				return err
			}

			fmt.Println(config.Passphrases[id])
			return nil
		},
	}

	root := &ffcli.Command{
		ShortUsage:  "pwinit <subcommand> [flags] [args...]",
		LongHelp:    "More info here: https://github.com/pathwar/pathwar/wiki/CLI#pwinit",
		Subcommands: []*ffcli.Command{entrypoint, env, config, passphrase},
		Exec: func(ctx context.Context, args []string) error {
			fmt.Println(banner.Inline("pwinit"))
			return flag.ErrHelp
		},
	}

	args := os.Args[1:]
	if len(args) > 0 && args[0] == "entrypoint" && len(args) > 1 && args[1] != "--" {
		args = append([]string{"entrypoint", "--"}, args[1:]...)
	}
	ctx := context.Background()
	if err := root.ParseAndRun(ctx, args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		log.Fatalf("fatal: %+v", err)
	}
}

func getConfig() (*pwinit.InitConfig, error) {
	configJSON, err := ioutil.ReadFile("/pwinit/config.json")
	if err != nil {
		return nil, err
	}
	var config pwinit.InitConfig
	if err := json.Unmarshal(configJSON, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
