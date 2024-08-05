// Copyright (C) 2024 Tianzhenxiong
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
