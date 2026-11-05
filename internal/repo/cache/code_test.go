package cache

import (
  "context"
  "errors"
  "testing"
  "weibook/internal/repo/cache/redismocks"

  "github.com/redis/go-redis/v9"
  "github.com/stretchr/testify/assert"
  "go.uber.org/mock/gomock"
)

func TestRedisCodeCache_Set(t *testing.T) {
  tcs := []struct {
    name string
    mock func(c *gomock.Controller) redis.Cmdable

    biz   string
    phone string
    code  string

    wantErr error
  }{
    {
      name: "成功",
      mock: func(c *gomock.Controller) redis.Cmdable {
        cmd := redismocks.NewMockCmdable(c)
        res := redis.NewCmd(context.Background())
        res.SetErr(nil)
        res.SetVal(int64(0))
        cmd.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:15992108773"}, "123123").
          Return(res)
        return cmd
      },

      biz:   "login",
      phone: "15992108773",
      code:  "123123",

      wantErr: nil,
    },
    {
      name: "redis错误",
      mock: func(c *gomock.Controller) redis.Cmdable {
        cmd := redismocks.NewMockCmdable(c)
        res := redis.NewCmd(context.Background())
        res.SetErr(errors.New("redis error"))
        //res.SetVal(int64())
        cmd.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:15992108773"}, "123123").
          Return(res)
        return cmd
      },

      biz:   "login",
      phone: "15992108773",
      code:  "123123",

      wantErr: errors.New("redis error"),
    },
    {
      name: "发送频繁",
      mock: func(c *gomock.Controller) redis.Cmdable {
        cmd := redismocks.NewMockCmdable(c)
        res := redis.NewCmd(context.Background())
        //res.SetErr(nil)
        res.SetVal(int64(-1))
        cmd.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:15992108773"}, "123123").
          Return(res)
        return cmd
      },

      biz:   "login",
      phone: "15992108773",
      code:  "123123",

      wantErr: ErrCodeSendTooMany,
    },
    {
      name: "系统错误",
      mock: func(c *gomock.Controller) redis.Cmdable {
        cmd := redismocks.NewMockCmdable(c)
        res := redis.NewCmd(context.Background())
        //res.SetErr(nil)
        res.SetVal(int64(-999))
        cmd.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:15992108773"}, "123123").
          Return(res)
        return cmd
      },

      biz:   "login",
      phone: "15992108773",
      code:  "123123",

      wantErr: errors.New("系统错误"),
    },
  }

  for _, tt := range tcs {
    t.Run(tt.name, func(t *testing.T) {
      c := gomock.NewController(t)
      defer c.Finish()

      cache := NewCodeCache(tt.mock(c))
      err := cache.Set(context.Background(), tt.biz, tt.phone, tt.code)
      assert.Equal(t, tt.wantErr, err)
    })
  }
}
