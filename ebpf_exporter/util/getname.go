package util

import (
	"context"

	"github.com/docker/docker/client"
)

func GetContainerName(cid string) string {
	// 创建 Docker 客户端
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "None"
	}
	cli.NegotiateAPIVersion(context.Background())

	// 获取容器详细信息
	containerJSON, err := cli.ContainerInspect(context.Background(), cid)
	if err != nil {
		return "None"
	}

	// 返回容器名称
	return containerJSON.Name
}
