package util

import (
	"os"

	"github.com/docker/docker/client"
)

var cli *client.Client

func init() {
	os.Setenv("DOCKER_API_VERSION", "1.45")
	cli, _ = client.NewClientWithOpts(client.FromEnv)
}
