package auth

import (
  "context"
  "errors"
  "weibook/internal/service/sms"

  "github.com/golang-jwt/jwt/v5"
)

type SMSService struct {
  svc sms.Service
  key string
}

// 权限校验，不同业务token不同(这里用biz作为token）
func (s *SMSService) Send(ctx context.Context, biz string, args []string, numbers ...string) error {

  var tc Claims
  token, err := jwt.ParseWithClaims(biz, &tc, func(token *jwt.Token) (any, error) {
    return s.key, nil
  })
  if err != nil {
    return err
  }
  if !token.Valid {
    return errors.New("token不合法")
  }

  return s.svc.Send(ctx, tc.Tpl, args, numbers...)
}

type Claims struct {
  jwt.RegisteredClaims
  Tpl string
}
