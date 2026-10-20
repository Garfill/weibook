package ioc

import (
  "weibook/internal/service/sms"
  "weibook/internal/service/sms/memory"
)

func InitSmsService() sms.Service {
  //return tencent.NewTencentService()
  return memory.NewMemorySms()
}
