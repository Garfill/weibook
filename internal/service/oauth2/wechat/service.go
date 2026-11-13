package serviceOauthWechat

import (
  "context"
  "errors"
  "fmt"
  "net/http"
  "net/url"
  "weibook/internal/domain"

  "github.com/goccy/go-json"
)

type Service interface {
  AuthURL(ctx context.Context) (string, error)
  VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error)
}

type WechatOauthService struct {
  appId string
}

func NewWechatOauthService(appId string) Service {
  return &WechatOauthService{
    appId: appId,
  }
}

func (w *WechatOauthService) AuthURL(ctx context.Context) (string, error) {
  const urlPattern = "https://open.wechat.com/qrconnect?appid=%s&redirect=%s"
  redirect := url.PathEscape("https://www.yourweb.com/callback")
  return fmt.Sprintf(urlPattern, w.appId, redirect), nil
}

func (w *WechatOauthService) VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error) {
  const targetPattern = "https://wechat.com/oauth2/token?=appId=%s&code=%s"
  var targetUrl = fmt.Sprintf(targetPattern, w.appId, code)
  resp, err := http.Get(targetUrl)
  defer resp.Body.Close()
  if err != nil {
    return domain.WechatInfo{}, err
  }

  var res WechatResult
  decoder := json.NewDecoder(resp.Body)
  err = decoder.Decode(&res)
  if err != nil {
    return domain.WechatInfo{}, err
  }

  if res.ErrCode != 0 {
    return domain.WechatInfo{}, errors.New("微信错误")
  }

  return domain.WechatInfo{}, nil

}

type WechatResult struct {
  ErrCode int64  `json:"errcode"`
  ErrMsg  string `json:"errmsg"`

  AccessToken  string `json:"access_token"`
  RefreshToken string `json:"refresh_token"`
  ExpireIn     string `json:"expire_in"`

  OpenId  string `json:"openid"`
  UnionId string `json:"unionid"`
}
