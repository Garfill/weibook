package middleware

import (
  "github.com/gin-contrib/sessions"
  "github.com/gin-gonic/gin"
  "net/http"
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
  return func(c *gin.Context) {
    reqPath := c.Request.URL.Path
    for _, p := range l.whitePaths {
      if reqPath == p {
        return
      }
    }
    sess := sessions.Default(c)
    id := sess.Get("userId")
    if id == nil {
      // session 内没有信息就是没有登录
      c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "登录失效，请重新登陆"})
      return
    }
  }
}
