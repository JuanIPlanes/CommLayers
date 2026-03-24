package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client  *redis.Client
	enabled bool
}

func New(url string) (*Cache, error) {
	if url == "" {
		return &Cache{}, nil
	}
	options, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(options)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &Cache{client: client, enabled: true}, nil
}

func (c *Cache) Enabled() bool {
	return c != nil && c.enabled
}

func (c *Cache) Save(ctx context.Context, key string, payload any, ttl time.Duration) error {
	if !c.Enabled() {
		return nil
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, body, ttl).Err()
}

func (c *Cache) Load(ctx context.Context, key string, dst any) (bool, error) {
	if !c.Enabled() {
		return false, nil
	}
	body, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	if err := json.Unmarshal(body, dst); err != nil {
		return false, err
	}
	return true, nil
}
