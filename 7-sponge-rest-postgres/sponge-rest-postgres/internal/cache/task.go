package cache

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/go-dev-frame/sponge/pkg/cache"
	"github.com/go-dev-frame/sponge/pkg/encoding"
	"github.com/go-dev-frame/sponge/pkg/utils"

	"sponge-rest-postgres/internal/database"
	"sponge-rest-postgres/internal/model"
)

const (
	// cache prefix key, must end with a colon
	taskCachePrefixKey = "task:"
	// TaskExpireTime expire time
	TaskExpireTime = 5 * time.Minute
)

var _ TaskCache = (*taskCache)(nil)

// TaskCache cache interface
type TaskCache interface {
	Set(ctx context.Context, id uint64, data *model.Task, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Task, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Task, error)
	MultiSet(ctx context.Context, data []*model.Task, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// taskCache define a cache struct
type taskCache struct {
	cache cache.Cache
}

// NewTaskCache new a cache
func NewTaskCache(cacheType *database.CacheType) TaskCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Task{}
		})
		return &taskCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Task{}
		})
		return &taskCache{cache: c}
	}

	return nil // no cache
}

// GetTaskCacheKey cache key
func (c *taskCache) GetTaskCacheKey(id uint64) string {
	return taskCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *taskCache) Set(ctx context.Context, id uint64, data *model.Task, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetTaskCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *taskCache) Get(ctx context.Context, id uint64) (*model.Task, error) {
	var data *model.Task
	cacheKey := c.GetTaskCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *taskCache) MultiSet(ctx context.Context, data []*model.Task, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetTaskCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *taskCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Task, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetTaskCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Task)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Task)
	for _, id := range ids {
		val, ok := itemMap[c.GetTaskCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *taskCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetTaskCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *taskCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetTaskCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *taskCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
