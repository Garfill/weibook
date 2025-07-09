package user

import (
  regexp "github.com/dlclark/regexp2"
  "github.com/gin-gonic/gin"
  "net/http"
  "weibook/internal/domain"
  "weibook/internal/service"
)

type UserHandler struct {
  svc            *service.UserService
  passwordRegexp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
  const passwordExp = `^(?=.*[a-z])(?=.*[A-Z])(?=.*[\d])[a-zA-Z\d]{8,}$`
  return &UserHandler{
    passwordRegexp: regexp.MustCompile(passwordExp, 0),
    svc:            svc,
  }
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
  group := server.Group("/user")
  group.POST("/login", u.Login)
  group.POST("/register", u.Register)
  group.POST("/logout", u.Logout)
  group.POST("/Edit", u.Edit)
  group.GET("/profile", u.GetProfile)
}

func (u *UserHandler) Register(ctx *gin.Context) {
  // 获取请求参数
  type ReqType struct {
    Name     string `json:"name"`
    Password string `json:"password"`
  }
  var req ReqType

  if err := ctx.Bind(&req); err != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": "系统解析参数错误"})
    return
  }

  // 校验密码包含大小写和数字
  ok, err := u.passwordRegexp.MatchString(req.Password)
  if err != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{"error": "内部正则错误"})
    return
  }
  if !ok {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": "请输入正确格式密码"})
    return
  }

  // 调用service
  err = u.svc.SignUp(ctx, domain.User{
    Name:     req.Name,
    Password: req.Password,
  })
  if err != nil {
    ctx.JSON(http.StatusOK, gin.H{"error": "系统错误"})
    return
  }

  ctx.JSON(http.StatusOK, gin.H{
    "msg": "register success",
  })
}

func (u *UserHandler) Login(ctx *gin.Context) {}

func (u *UserHandler) Logout(ctx *gin.Context) {}

func (u *UserHandler) Edit(ctx *gin.Context) {}

func (u *UserHandler) GetProfile(ctx *gin.Context) {
  type User struct {
    Name string `json:"name"`
  }
  user := User{
    Name: "guan",
  }
  ctx.JSON(http.StatusOK, gin.H{
    "user": user,
  })
}
