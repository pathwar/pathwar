package pwhypervisor

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
	"go.uber.org/zap"
	"moul.io/godev"
	"pathwar.land/go/pkg/errcode"
	"pathwar.land/go/pkg/pwcompose"
)

func Daemon(ctx context.Context, cli *client.Client, logger *zap.Logger) error {
	info, err := pwcompose.GetPathwarInfo(ctx, cli)
	if err != nil {
		return errcode.ErrComposeGetPathwarInfo.Wrap(err)
	}
	fmt.Println("info", godev.PrettyJSON(info))
	return errcode.ErrNotImplemented
}
