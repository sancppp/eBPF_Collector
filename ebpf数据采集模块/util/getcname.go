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
