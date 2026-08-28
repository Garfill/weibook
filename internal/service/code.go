package service

import (
  "context"
  "fmt"
  "math/rand"
  "weibook/internal/repo"
  "weibook/internal/service/sms"
)

// 假设验证码存在 10 分钟 有效期，存 redis
// redisKey = phone_code:$biz:$phone
// 1. redis 没有这个key，直接发送
// 2. 有key 没有过期时间，系统有异常
// 3. 有key，还没过期(过期时间大于9分钟），发送频繁，拒绝发送
// 4. 重新发送

// check-do-sth 模式下注意分布式环境下并发安全问题
// 多线程同时查询相同号码，会有并发问题，两个都返回没发过
// 非语法层面问题，不能使用 channel或sync.lock

// 解决
// 1. 分布式锁：性能差，难用
// 2. redis内使用lua脚本：redis是单线程

type CodeService struct {
  repo repo.CodeRepo
  sms  sms.Service
}

// 发验证码
func (svc *CodeService) Send(
  ctx context.Context,
  biz string, // 不同业务场景 biz 不一样（登录/更改密码等）
  phone string,
) error {
  //  1. 生成验证码
  code := svc.generateCode()
  //  2. 存储到redis
  err := svc.repo.Store(ctx, biz, phone, code)
  if err != nil {
    // 有问题
    return err
  }
  //  3.发送
  err = svc.sms.Send(ctx, "1812222", []string{code}, phone)
  //if err != nil {
  //  //  redis有验证码
  //  //  1.可能超时error，不明确是否失败
  //  //  2.如果重试，最好在初始化时候传入 smsRetryService
  //}
  return err
}

// 校验验证码
func (svc *CodeService) Verify(
  ctx context.Context,
  biz string, // 不同业务场景 biz 不一样（登录/更改密码等）
  code string,
  phone string,
) (bool, error) {
  return svc.repo.Verify(ctx, biz, phone, code)
}

func (svc *CodeService) generateCode() string {
  // 0 - 999999随机
  num := rand.Intn(1000000)
  // 不够6位数的前面补充0
  return fmt.Sprintf("%6d", num)
}
