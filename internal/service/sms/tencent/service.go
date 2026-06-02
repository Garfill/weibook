package tencent

import (
  "context"
  "fmt"

  tSms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
  appId    *string
  signName *string
  client   *tSms.Client
}

func (s *Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
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

func NewService(client *tSms.Client, appId string, signName string) *Service {
  return &Service{
    appId:    &appId,
    signName: &signName,
    client:   client,
  }
}
