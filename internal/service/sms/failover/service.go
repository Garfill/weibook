package failover

import (
  "context"
  "errors"
  "sync/atomic"
  "weibook/internal/service/sms"
)

// 装饰器
// 服务商失败处理
type FailoverSMSService struct {
  svcs []sms.Service
  idx  int64
}

func (f *FailoverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
  // 循环发送每个服务商
  for _, svc := range f.svcs {
    err := svc.Send(ctx, tpl, args, numbers...)
    // 发送成功
    if err == nil {
      return nil
    }
    // 发送失败，日志监控
  }
  return errors.New("全部服务商失败")
}

func (f *FailoverSMSService) SendV1(ctx context.Context, tpl string, args []string, numbers ...string) error {
  // 循环发送每个服务商
  idx := atomic.AddInt64(&f.idx, 1)
  length := int64(len(f.svcs))
  // 加个 idx 控制发送位置
  for i := idx; i < length+idx; i++ {
    err := f.svcs[i%length].Send(ctx, tpl, args, numbers...)
    // 发送成功
    switch err {
    case nil:
      return nil
    case context.Canceled, context.DeadlineExceeded:
      return err
    }
    // 发送失败，日志监控
    // 下一轮循环
  }
  return errors.New("全部服务商失败")
}

func NewFailoverSMSService(svcs []sms.Service) sms.Service {
  return &FailoverSMSService{
    svcs: svcs,
  }
}
