package util

import (
	"context"

	"github.com/docker/docker/api/types/container"
)

// ContainerInfo holds the information of a container
type ContainerInfo struct {
	Name string
	CID  string
	IP   string
}

// GetContainerInfo returns a slice of ContainerInfo for all running containers
func GetContainerInfo() ([]ContainerInfo, error) {

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		return nil, err
	}

	var containerInfos []ContainerInfo

	for _, container := range containers {
		containerJSON, err := cli.ContainerInspect(context.Background(), container.ID)
		if err != nil {
			return nil, err
		}

		ip := ""
		if containerJSON.NetworkSettings != nil {
			for _, network := range containerJSON.NetworkSettings.Networks {
				ip = network.IPAddress
				break
			}
		}

		containerInfos = append(containerInfos, ContainerInfo{
			Name: container.Names[0],
			CID:  container.ID,
			IP:   ip,
		})
	}

	return containerInfos, nil
}
