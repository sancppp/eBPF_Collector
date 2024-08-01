package util

import (
	"context"
	"sync"

	"github.com/docker/docker/client"
)

var once sync.Once
var cli *client.Client

func GetContainerName(cid string) string {
	//只初始化一次cli
	once.Do(func() {
		// 创建 Docker 客户端
		cli, _ = client.NewClientWithOpts(client.FromEnv)
	})

	cli.NegotiateAPIVersion(context.Background())
	// 获取容器详细信息
	containerJSON, err := cli.ContainerInspect(context.Background(), cid)
	if err != nil {
		return "None"
	}
	// 返回容器名称
	return containerJSON.Name
}
