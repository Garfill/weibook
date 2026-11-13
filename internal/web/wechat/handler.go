package wechat

import (
  "net/http"
  "weibook/internal/service"
  "weibook/internal/service/oauth2/wechat"
  "weibook/internal/web"

  "github.com/gin-gonic/gin"
)

type OAuth2WechatHandler struct {
  wechatOAuthService serviceOauthWechat.Service
  userSvc            service.UserService
  web.JwtHandler
}

func NewOAuth2WechatHandler(svc serviceOauthWechat.Service, userSvc service.UserService) *OAuth2WechatHandler {
  return &OAuth2WechatHandler{
    wechatOAuthService: svc,
    userSvc:            userSvc,
  }
}

func (h *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
  g := server.Group("/oauth2/wechat")
  g.GET("/authurl", h.AuthURL)
  g.Any("/callback", h.Callback)
}

func (h *OAuth2WechatHandler) AuthURL(c *gin.Context) {
  url, err := h.wechatOAuthService.AuthURL(c)
  if err != nil {
    c.JSON(http.StatusOK, web.Result{
      Code: 400,
      Msg:  "生成url错误",
    })
    return
  }
  c.JSON(http.StatusOK, web.Result{
    Code: 200,
    Data: url,
  })
}

// 微信回调
func (h *OAuth2WechatHandler) Callback(ctx *gin.Context) {
  code := ctx.Query("code")
  info, err := h.wechatOAuthService.VerifyCode(ctx, code)
  if err != nil {
    ctx.JSON(http.StatusOK, web.Result{
      Code: 500,
      Msg:  "系统错误",
    })
    return
  }

  user, err := h.userSvc.FindOrCreateByWechat(ctx, info)
  if err != nil {
    ctx.JSON(http.StatusOK, web.Result{
      Code: 500,
      Msg:  "系统错误",
    })
    return
  }

  err = h.JwtHandler.SetJwtToken(ctx, user.Id)
  if err != nil {
    ctx.JSON(http.StatusOK, web.Result{
      Code: 500,
      Msg:  "系统错误",
    })
    return
  }

  ctx.JSON(http.StatusOK, web.Result{
    Code: 200,
  })

}
