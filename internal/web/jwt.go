package web

import (
  "time"
  "weibook/internal/variable"

  "github.com/gin-gonic/gin"
  "github.com/golang-jwt/jwt/v5"
)

type JwtHandler struct {
}

func (h JwtHandler) SetJwtToken(ctx *gin.Context, uid int64) error {
  userClaim := UserClaim{
    Uid: uid,
    RegisteredClaims: jwt.RegisteredClaims{
      // jwt 过期时间
      // 60m 过期
      ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
    },
    UserAgent: ctx.Request.UserAgent(), // ua 用来判断是否异地登录
  }

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim)
  // 2. 加签，返回token字符
  tokenStr, err := token.SignedString(variable.JWTEncryptKey)
  if err != nil {
    return err
  }
  ctx.Header("x-jwt-token", tokenStr) // 返回前端存储，下次请求通过请求头带上来
  return nil
}

// jwt 中 用户信息的 声明
type UserClaim struct {
  jwt.RegisteredClaims
  Uid       int64
  UserAgent string
}
