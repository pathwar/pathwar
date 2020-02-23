package main

import (
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

	"github.com/peterbourgon/ff/ffcli"
	"pathwar.land/v2/go/pkg/errcode"
	"pathwar.land/v2/go/pkg/pwinit"
)

func main() {
	log.SetFlags(0)

	entrypoint := &ffcli.Command{
		Name:  "entrypoint",
		Usage: "pwinit entrypoint [args...]",
		Exec: func(args []string) error {
			// FIXME: lock to block other commands

			_, err := os.Stat("/pwinit/config.json")
			if err != nil {
				log.Printf("no such config file, skipping on-init hook (%v)", err)
			} else {
				log.Print("starting on-init hook")
				// prepare the challenge
				err = os.Chmod("/pwinit/on-init", 0555)
				if err != nil {
					return errcode.ErrExecuteOnInitHook.Wrap(err)
				}
				cmd := exec.Command("/pwinit/on-init")
				err = cmd.Run()
				if err != nil {
					return errcode.ErrExecuteOnInitHook.Wrap(err)
				}

				// clean pwinit config file that contains passphrases
				for _, filename := range []string{"/pwinit/config.json", "/pwinit/on-init"} {
					err = os.Remove(filename)
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
		Name:  "env",
		Usage: "pwinit entrypoint [args...]",
		Exec: func([]string) error {
			for _, line := range os.Environ() {
				fmt.Println(line)
			}
			return nil
		},
	}

	config := &ffcli.Command{
		Name:  "config",
		Usage: "pwinit config [args...]",
		Exec: func([]string) error {
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
		Name:  "passphrase",
		Usage: "pwinit passphrase ID",
		Exec: func(args []string) error {
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
		Usage:       "pwinit <subcommand> [flags] [args...]",
		LongHelp:    "More info here: https://github.com/pathwar/pathwar/wiki/CLI#pwinit",
		Subcommands: []*ffcli.Command{entrypoint, env, config, passphrase},
		Exec:        func([]string) error { return flag.ErrHelp },
	}

	args := os.Args[1:]
	if args[0] == "entrypoint" && len(args) > 1 && args[1] != "--" {
		args = append([]string{"entrypoint", "--"}, args[1:]...)
	}
	if err := root.Run(args); err != nil {
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
