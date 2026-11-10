package failover

import (
  "context"
  "sync/atomic"
  "weibook/internal/service/sms"
)

type TimeOutFailoverSMSService struct {
  svcs []sms.Service

  // 当前服务序号标识
  idx int32
  // 连续超时个数
  cnt int32
  // 连续超时阈值
  threshold int32
}

func NewTimeoutFailoverSMSService(svcs []sms.Service, threshold int32) sms.Service {
  return &TimeOutFailoverSMSService{
    svcs:      svcs,
    threshold: threshold,
  }
}

func (t *TimeOutFailoverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
  cnt := atomic.LoadInt32(&t.cnt)
  idx := atomic.LoadInt32(&t.idx)

  if cnt >= t.threshold {
    // 当前服务超出阈值,切换服务
    newIdx := (idx + 1) % int32(len(t.svcs))
    // 这样原子操作，下一个服务进来才是用新的
    if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
      // 切换成功，重置
      atomic.StoreInt32(&t.cnt, 0)
    }
    idx = newIdx
  }

  svc := t.svcs[idx]
  err := svc.Send(ctx, tpl, args, numbers...)
  switch err {
  case nil:
    atomic.StoreInt32(&t.cnt, 0)
    return nil
  case context.DeadlineExceeded:
    atomic.AddInt32(&t.cnt, 1)
    return err
  default:
    // 非超时错误
    // 可增加cnt或者不增加
  }
  return err
}
