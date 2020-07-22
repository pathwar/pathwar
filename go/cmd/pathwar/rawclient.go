package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"

	"github.com/peterbourgon/ff"
	"github.com/peterbourgon/ff/ffcli"
	"moul.io/godev"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func rawclientCommand() *ffcli.Command {
	rawclientFlags := flag.NewFlagSet("client", flag.ExitOnError)
	rawclientFlags.StringVar(&httpAPIAddr, "http-api-addr", defaultHTTPApiAddr, "HTTP API address")
	rawclientFlags.StringVar(&ssoOpts.ClientID, "sso-clientid", ssoOpts.ClientID, "SSO ClientID")
	rawclientFlags.StringVar(&ssoOpts.ClientSecret, "sso-clientsecret", ssoOpts.ClientSecret, "SSO ClientSecret")
	rawclientFlags.StringVar(&ssoOpts.Realm, "sso-realm", ssoOpts.Realm, "SSO Realm")
	rawclientFlags.StringVar(&ssoOpts.TokenFile, "sso-token-file", ssoOpts.TokenFile, "Token file")

	return &ffcli.Command{
		Name:      "rawclient",
		Usage:     "pathwar [global flags] rawclient [rawclient flags] <method> <path> [INPUT (json)]",
		ShortHelp: "make API calls",
		LongHelp: `EXAMPLES
  pathwar rawclient GET /user/session
  season=$(pathwar rawclient GET /user/session | jq -r '.seasons[0].season.id')
  pathwar rawclient GET "/season-challenges?season_id=$season"`,
		FlagSet: rawclientFlags,
		Options: []ff.Option{ff.WithEnvVarNoPrefix()},
		Exec: func(args []string) error {
			if len(args) < 2 || len(args) > 3 {
				return flag.ErrHelp
			}
			if err := globalPreRun(); err != nil {
				return err
			}
			ctx := context.Background()
			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			method := args[0]
			path := args[1]
			var input []byte
			if len(args) > 2 {
				input = []byte(args[2])
			}
			output, err := apiClient.Raw(ctx, method, path, input)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			var data interface{}
			err = json.Unmarshal(output, &data)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			fmt.Println(godev.PrettyJSON(data))
			return nil
		},
	}
}
