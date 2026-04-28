package cache

import (
  "context"
  "fmt"
  "time"
  "weibook/internal/domain"

  "github.com/goccy/go-json"
  "github.com/redis/go-redis/v9"
)

var ErrUserNotFound = redis.Nil // redis 特定错误

// 中间件团队最佳实现
// 统一内存，可以传入多个实现（redis，memcache, memory)

type MiddleCache interface {
  Get(ctx context.Context, key string) (any, error)
  Set(ctx context.Context, key string, val interface{}, expiration time.Duration)
}

type UserMiddleCache struct {
  cache MiddleCache
}

// ------end-----

// 较弱实现

type UserCache struct {
  // 面向redis接口，不区分redis类型（单机/集群）
  client     redis.Cmdable
  expiration time.Duration // cache 过期时间
}

func NewUserCache(client redis.Cmdable) *UserCache {
  // 外部注入 redis 客户端连接
  return &UserCache{
    client:     client,
    expiration: time.Minute * 10,
  }
}

// 如果有数据，err是nil
// 如果没数据，返回特定 err
func (cache *UserCache) Get(ctx context.Context, id int64) (domain.User, error) {
  key := cache.Key(id)
  // 数据不存在，err == ErrUserNotFound，返回上级
  val, err := cache.client.Get(ctx, key).Bytes()
  if err != nil {
    return domain.User{}, err
  }
  var u domain.User
  err = json.Unmarshal(val, &u)
  return u, err
}

func (cache *UserCache) Set(ctx context.Context, u domain.User) error {
  val, err := json.Marshal(u)
  if err != nil {
    return err
  }
  key := cache.Key(u.Id)
  return cache.client.Set(ctx, key, val, cache.expiration).Err()
}

func (cache *UserCache) Key(id int64) string {
  return fmt.Sprint("user:info:", id)
}
