package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/syumai/workers/cloudflare/kv"
)

// KVClient Cloudflare KV 存储客户端
type KVClient struct {
	namespace *kv.Namespace
}

// NewKVClient 创建KV客户端实例
func NewKVClient(bindingName string) (*KVClient, error) {
	// 从Worker环境中获取KV命名空间绑定
	namespace, err := kv.NewNamespace(bindingName)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to KV namespace: %w", err)
	}

	return &KVClient{namespace: namespace}, nil
}

// Get 获取KV存储的值
func (c *KVClient) Get(ctx context.Context, key string) (string, error) {
	value, err := c.namespace.GetString(key, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get key %s: %w", key, err)
	}
	return value, nil
}

// GetJSON 获取KV存储的JSON对象并反序列化
func (c *KVClient) GetJSON(ctx context.Context, key string, v interface{}) error {
	value, err := c.Get(ctx, key)
	if err != nil {
		return err
	}

	if value == "" {
		return fmt.Errorf("key %s not found", key)
	}

	err = json.Unmarshal([]byte(value), v)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON for key %s: %w", key, err)
	}

	return nil
}

// Put 存储值到KV
// expirationTTL: 过期时间（秒），0表示永不过期
func (c *KVClient) Put(ctx context.Context, key string, value string, expirationTTL int) error {
	var options *kv.PutOptions
	if expirationTTL > 0 {
		options = &kv.PutOptions{
			ExpirationTTL: expirationTTL,
		}
	}

	err := c.namespace.PutString(key, value, options)
	if err != nil {
		return fmt.Errorf("failed to put key %s: %w", key, err)
	}

	return nil
}

// PutJSON 将对象序列化为JSON后存储到KV
func (c *KVClient) PutJSON(ctx context.Context, key string, v interface{}, expirationTTL int) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON for key %s: %w", key, err)
	}

	return c.Put(ctx, key, string(bytes), expirationTTL)
}

// Delete 删除KV中的键
func (c *KVClient) Delete(ctx context.Context, key string) error {
	err := c.namespace.Delete(key)
	if err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}
	return nil
}

// Exists 检查键是否存在
func (c *KVClient) Exists(ctx context.Context, key string) (bool, error) {
	value, err := c.Get(ctx, key)
	if err != nil {
		// 如果是"key not found"错误，返回false
		return false, nil
	}
	return value != "", nil
}

// CacheKey 生成缓存键名
func CacheKey(prefix string, id string) string {
	return fmt.Sprintf("%s:%s", prefix, id)
}

// CacheKeyWithTimestamp 生成带时间戳的缓存键名（用于防止缓存穿透）
func CacheKeyWithTimestamp(prefix string, id string) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s:%s:%d", prefix, id, timestamp)
}

// 常用缓存键前缀
const (
	CachePrefixSeiyuu  = "seiyuu"
	CachePrefixGroup   = "group"
	CachePrefixMoegirl = "moegirl"
)

// 默认缓存过期时间（秒）
const (
	CacheTTLShort  = 300    // 5分钟
	CacheTTLMedium = 1800   // 30分钟
	CacheTTLLong   = 3600   // 1小时
	CacheTTLDay    = 86400  // 1天
	CacheTTLWeek   = 604800 // 1周
)
