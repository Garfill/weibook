package cache

import (
  "context"
  _ "embed"
  "errors"
  "fmt"

  "github.com/redis/go-redis/v9"
)

var (
  ErrCodeSendTooMany   = errors.New("发送太频繁")
  ErrCodeVerifyTooMany = errors.New("验证太频繁")
)

// 编译器在编译时，会将代码放到luaCode变量
//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

type CodeCache interface {
  Set(ctx context.Context, biz, phone, code string) error
  Verify(ctx context.Context, biz, phone, code string) (bool, error)
  Key(biz, phone string) string
}

func NewCodeCache(client redis.Cmdable) CodeCache {
  return &RedisCodeCache{
    client: client,
  }
}

type RedisCodeCache struct {
  client redis.Cmdable
}

func (c *RedisCodeCache) Set(ctx context.Context, biz, phone, code string) error {
  res, err := c.client.Eval(ctx, luaSetCode, []string{c.Key(biz, phone)}, code).Int() // lua脚本返回的就是int
  if err != nil {
    return err
  }
  switch res {
  case 0:
    //  没问题
    return nil
  case -1:
    //  频繁
    return ErrCodeSendTooMany
  default:
    // 其他都可以认为系统错误
    return errors.New("系统错误")
  }
}

func (c *RedisCodeCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
  res, err := c.client.Eval(ctx, luaVerifyCode, []string{c.Key(biz, phone)}, code).Int()
  if err != nil {
    return false, err
  }
  switch res {
  case 0:
    return true, nil
  case -1:
    //  一直输入错误
    //  清空校验key，次数
    return false, ErrCodeVerifyTooMany
  case -2:
    // 减少一次校验次数
    return false, nil
  default:
    return false, nil
  }
}

func (c *RedisCodeCache) Key(biz, phone string) string {
  return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
