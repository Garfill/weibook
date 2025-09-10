package middleware

import (
  "fmt"
  "net/http"
  "weibook/internal/variable"

  "github.com/gin-gonic/gin"
  "github.com/golang-jwt/jwt/v5"
)

type LoginJWTMiddleWareBuilder struct {
  whitePaths []string
}

func NewLoginJWTMiddleBuilder() *LoginJWTMiddleWareBuilder {
  return &LoginJWTMiddleWareBuilder{
    whitePaths: []string{"/user/login", "/user/register"},
  }
}

// 返回 l 是为链式调用
func (l *LoginJWTMiddleWareBuilder) IgnorePath(p string) *LoginJWTMiddleWareBuilder {
  l.whitePaths = append(l.whitePaths, p)
  return l
}

func (l *LoginJWTMiddleWareBuilder) Build() gin.HandlerFunc {
  return func(c *gin.Context) {
    reqPath := c.Request.URL.Path
    for _, p := range l.whitePaths {
      if reqPath == p {
        return
      }
    }
    // 使用jwt 校验
    tokenStr := c.GetHeader("x-jwt-token")
    if tokenStr == "" {
      // 没token 就是没登录
      c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "登录失效"})
      return
    }
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
      return variable.JWTEncryptKey, nil
    })
    if err != nil {
      c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "登录失效"})
      return
    }
    if token == nil || !token.Valid {
      c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "登录失效"})
      return
    }
    fmt.Println("token ======", tokenStr)
  }
}
