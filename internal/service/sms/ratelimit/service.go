package ratelimit

import (
  "context"
  "errors"
  "weibook/internal/rateLimit"
  "weibook/internal/service/sms"
)

type RateLimitSmsService struct {
  smsSvc  sms.Service
  limiter rateLimit.Limiter
  key     string
}

func NewRateLimitSmsService(s sms.Service, l rateLimit.Limiter) sms.Service {
  return &RateLimitSmsService{
    smsSvc:  s,
    limiter: l,
    key:     "sms-limiter",
  }
}

var ErrSMSRateLimit = errors.New("短信触发限流")

// 实现了 Service 抽象接口
// 装饰器写法
// 在核心实现外面包一层实现
func (s *RateLimitSmsService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
  // 追加新特性代码
  limited, err := s.limiter.Limited(ctx, s.key)
  if err != nil {
    // 系统错误
    // 可以限流：保守策略，下游系统不给力
    //  不限流：容错，业务可用性要求很高的时候不限流
    return err
  }
  if limited {
    // 限流了
    return ErrSMSRateLimit
  }
  // 核心实现
  err = s.smsSvc.Send(ctx, tpl, args, numbers...)
  return err
}
