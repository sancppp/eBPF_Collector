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
	"sync"
	"time"
)

// 定义一个结构体来存储缓存条目
type cacheEntry struct {
	name      string
	timestamp time.Time
}

// 创建一个全局缓存和互斥锁
var (
	cache      = make(map[string]cacheEntry)
	cacheMutex = &sync.Mutex{}
)

// 缓存有效期为2小时
const cacheDuration = 2 * time.Hour

// 获取容器名称
func GetContainerName(cid string) string {
	cacheMutex.Lock()
	entry, found := cache[cid]
	cacheMutex.Unlock()

	// 检查缓存是否存在且未过期
	if found && time.Since(entry.timestamp) < cacheDuration {
		return entry.name
	}

	cli.NegotiateAPIVersion(context.Background())
	// 获取容器详细信息
	containerJSON, err := cli.ContainerInspect(context.Background(), cid)
	if err != nil {
		return "None"
	}

	// 更新缓存
	cacheMutex.Lock()
	cache[cid] = cacheEntry{name: containerJSON.Name, timestamp: time.Now()}
	cacheMutex.Unlock()

	// 返回容器名称
	return containerJSON.Name
}
