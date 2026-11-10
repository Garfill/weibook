package tencent

import (
  "context"
  "fmt"
  "weibook/internal/rateLimit"

  tSms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
  appId    *string
  signName *string
  client   *tSms.Client
  limiter  rateLimit.Limiter
}

// 后期通过 test 单元测试
func (s *Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
  limited, err := s.limiter.Limited(ctx, "tencent:sms")
  if err != nil {
    // 系统错误
    // 可以限流：保守策略，下游系统不给力
    //  不限流：容错，业务可用性要求很高的时候不限流
    return fmt.Errorf("短信限流错误，%w", err)
  }
  if limited {
    // 限流了
    return fmt.Errorf("触发限流")
  }
  req := tSms.NewSendSmsRequest()
  req.SmsSdkAppId = s.appId
  req.SignName = s.signName
  req.TemplateId = &tpl
  req.PhoneNumberSet = s.toStringPtrSlice(numbers)
  req.TemplateParamSet = s.toStringPtrSlice(args)
  resp, err := s.client.SendSms(req)
  if err != nil {
    return err
  }
  for _, status := range resp.Response.SendStatusSet {
    if status.Code == nil || *(status.Code) != "Ok" {
      return fmt.Errorf("发送短信失败 %s %s", *status.Code, *status.Message)
    }
  }
  return nil
}

func (s *Service) toStringPtrSlice(src []string) []*string {
  ptr := make([]*string, len(src))
  for i := range src {
    ptr = append(ptr, &src[i])
  }
  return ptr
}

func NewTencentService(client *tSms.Client, appId string, signName string) *Service {
  return &Service{
    appId:    &appId,
    signName: &signName,
    client:   client,
  }
}
