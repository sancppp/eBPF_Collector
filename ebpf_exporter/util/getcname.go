package util

import (
	"context"
)

func GetContainerName(cid string) string {

	cli.NegotiateAPIVersion(context.Background())
	// 获取容器详细信息
	containerJSON, err := cli.ContainerInspect(context.Background(), cid)
	if err != nil {
		return "None"
	}
	// 返回容器名称
	return containerJSON.Name
}
