package middleware

import (
  "net/http"
  "time"
  "weibook/internal/variable"
  "weibook/internal/www/user"

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
        println("reqPath === ", reqPath)
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
    var userClaim user.UserClaim
    token, err := jwt.ParseWithClaims(tokenStr, &userClaim, func(token *jwt.Token) (any, error) {
      return variable.JWTEncryptKey, nil
    })
    if err != nil {
      c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "登录失效"})
      return
    }
    if userClaim.UserAgent != c.Request.UserAgent() {
      // UA 变化，可能token泄漏
      // 日志记录
      c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "登录失效"})
      return
    }
    // valid 会校验过期时间
    if token == nil || !token.Valid || userClaim.Uid == 0 {
      c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "登录失效"})
      return
    }

    // 刷新 jwt 过期时间
    now := time.Now()
    if userClaim.ExpiresAt.Sub(now) < time.Second*50 {
      // 每 10s 刷新
      userClaim.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * 1))

      // 生成新 token
      newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim)
      newTokenStr, err := newToken.SignedString(variable.JWTEncryptKey)
      if err != nil {
        // 不return，日志记录失败
        c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "token续约失败"})
      }
      c.Header("x-jwt-token", newTokenStr)
    }

    // 拿到 token 内的信息 claim，通过 context 传递到接口上
    c.Set("userInfo", userClaim)
  }
}
