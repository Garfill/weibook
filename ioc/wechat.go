package ioc

import (
  "os"
  serviceOauthWechat "weibook/internal/service/oauth2/wechat"
)

func InitWechatOAuth2() serviceOauthWechat.Service {
  appId, ok := os.LookupEnv("WECHAT_APP_ID")
  if !ok {
    panic("读取失败")
  }
  return serviceOauthWechat.NewWechatOauthService(appId)
}
