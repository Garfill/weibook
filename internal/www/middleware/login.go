package middleware

import (
  "encoding/gob"
  "net/http"
  "time"

  "github.com/gin-contrib/sessions"
  "github.com/gin-gonic/gin"
)

type LoginMiddleWareBuilder struct {
  whitePaths []string
}

func NewLoginMiddleBuilder() *LoginMiddleWareBuilder {
  return &LoginMiddleWareBuilder{
    whitePaths: []string{"/user/login", "/user/register"},
  }
}

// 返回 l 是为链式调用
func (l *LoginMiddleWareBuilder) IgnorePath(p string) *LoginMiddleWareBuilder {
  l.whitePaths = append(l.whitePaths, p)
  return l
}

func (l *LoginMiddleWareBuilder) Build() gin.HandlerFunc {
  // 注册这个类型用于序列化
  gob.Register(time.Now())
  return func(c *gin.Context) {
    reqPath := c.Request.URL.Path
    for _, p := range l.whitePaths {
      if reqPath == p {
        return
      }
    }
    // session内会存储用户登录相关的信息
    sess := sessions.Default(c)
    // session 内没有信息就是没有登录
    uid := sess.Get("userId")
    if uid == nil {
      c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "登录失效，请重新登陆"})
      return
    }
    sess.Set("userId", uid)
    sess.Options(sessions.Options{
      MaxAge:   60,
      Path:     "/",
      HttpOnly: true,
    })
    // session 存储cookie刷新时间，也可以是token刷新时间，防止登录频繁过期
    const updateTimeKey = "update_time"
    now := time.Now().UnixMilli()
    lastUpdateTime := sess.Get(updateTimeKey)
    if lastUpdateTime == nil {
      // 刚登录
      sess.Set(updateTimeKey, now)
      if err := sess.Save(); err != nil {
        panic(err)
      }
      return
    }
    lastUpdateTimeValue, ok := lastUpdateTime.(int64)
    if !ok {
      c.AbortWithStatus(500)
      return
    }
    if now-lastUpdateTimeValue > 1000*5 {
      sess.Set(updateTimeKey, now)
      err := sess.Save()
      if err != nil {
        panic(err)
      }
      return
    }
  }
}
