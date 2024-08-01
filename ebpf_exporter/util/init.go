package util

import (
	"github.com/docker/docker/client"
)

var cli *client.Client

func init() {
	cli, _ = client.NewClientWithOpts(client.FromEnv)
}
