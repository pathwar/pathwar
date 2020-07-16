package pwcompose

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/client"
	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"go.uber.org/zap"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func PS(ctx context.Context, depth int, cli *client.Client, logger *zap.Logger) error {
	logger.Debug("ps", zap.Int("depth", depth))

	containersInfo, err := GetContainersInfo(ctx, cli)
	if err != nil {
		return errcode.ErrComposeGetContainersInfo.Wrap(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NGINX ID", "STATUS", "CREATED"})
	if containersInfo.NginxContainer.ID != "" {
		table.Append([]string{
			containersInfo.NginxContainer.ID[:7],
			strings.Replace(containersInfo.NginxContainer.Status, "Up ", "", 1),
			strings.Replace(humanize.Time(time.Unix(containersInfo.NginxContainer.Created, 0)), " ago", "", 1),
		})
	}
	table.Render()

	table = tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "CHALLENGE", "SVC", "PORTS", "STATUS", "CREATED"})

	for _, flavor := range containersInfo.RunningFlavors {
		for uid, container := range flavor.Containers {
			ports := []string{}
			for _, port := range container.Ports {
				if port.PublicPort != 0 {
					ports = append(ports, strconv.Itoa(int(port.PublicPort)))
				}
			}

			table.Append([]string{
				uid[:7],
				flavor.ChallengeID(),
				container.Labels[serviceNameLabel],
				strings.Join(ports, ", "),
				strings.Replace(container.Status, "Up ", "", 1),
				strings.Replace(humanize.Time(time.Unix(container.Created, 0)), " ago", "", 1),
			})
		}
	}
	table.Render()
	return nil
}
