package user

import (
  "errors"
  regexp "github.com/dlclark/regexp2"
  "github.com/gin-contrib/sessions"
  "github.com/gin-gonic/gin"
  "net/http"
  "time"
  "weibook/internal/domain"
  "weibook/internal/service"
)

type UserHandler struct {
  svc            *service.UserService
  passwordRegexp *regexp.Regexp
  emailRegexp    *regexp.Regexp
}

var ErrDuplicateUser = service.ErrDuplicateUser

func NewUserHandler(svc *service.UserService) *UserHandler {
  // 限制 8-50 长度的密码，防止加密算法不支持
  const passwordExp = `^(?=.*[a-z])(?=.*[A-Z])(?=.*[\d])[~!@#$%^&a-zA-Z\d]{8,50}$`
  const emailExp = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
  return &UserHandler{
    passwordRegexp: regexp.MustCompile(passwordExp, 0),
    emailRegexp:    regexp.MustCompile(emailExp, 0),
    svc:            svc,
  }
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
  group := server.Group("/user")
  group.POST("/register", u.Register)
  group.POST("/login", u.Login)
  group.POST("/logout", u.Logout)
  group.POST("/Edit", u.Edit)
  group.GET("/profile", u.GetProfile)
}

func (u *UserHandler) Register(ctx *gin.Context) {
  // 获取请求参数
  type ReqType struct {
    Name     string `json:"name"`
    Password string `json:"password"`
    Email    string `json:"email"`
  }
  var req ReqType

  if err := ctx.Bind(&req); err != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
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

  // 校验邮箱
  ok, err = u.emailRegexp.MatchString(req.Email)
  if err != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{"error": "内部正则错误"})
    return
  }
  if !ok {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": "请输入正确格式的邮箱"})
    return
  }

  // 调用service
  err = u.svc.SignUp(ctx, domain.User{
    Name:     req.Name,
    Password: req.Password,
    Email:    req.Email,
  })

  // 返回的特定错误
  if errors.Is(err, ErrDuplicateUser) {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }

  if err != nil {
    ctx.JSON(http.StatusOK, gin.H{"error": "系统错误"})
    return
  }

  ctx.JSON(http.StatusOK, gin.H{
    "msg": "register success",
  })
}

func (u *UserHandler) Login(ctx *gin.Context) {
  type LoginReq struct {
    Email    string `json:"email"`
    Password string `json:"password"`
  }
  var req LoginReq
  err := ctx.Bind(&req)
  if err != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
    return
  }
  user, err := u.svc.Login(ctx, req.Email, req.Password)
  if errors.Is(err, service.ErrInvalidUserOrPwd) {
    ctx.JSON(http.StatusNotFound, gin.H{"error": "帐号或者密码错误"})
    return
  }
  if err != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误"})
    return
  }
  // 登录成功
  // 取出session并设置
  session := sessions.Default(ctx)
  session.Set("userId", user.Id)
  session.Set("refreshTime", time.Now().UnixMilli())
  session.Save()
  ctx.JSON(http.StatusOK, gin.H{"msg": "登录成功"})
  return
}

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
