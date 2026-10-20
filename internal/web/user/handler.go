package web

import (
  "errors"
  "fmt"
  "net/http"
  "strconv"
  "time"
  "weibook/internal/domain"
  "weibook/internal/service"
  "weibook/internal/variable"
  "weibook/internal/web"

  regexp "github.com/dlclark/regexp2"
  "github.com/gin-contrib/sessions"
  "github.com/gin-gonic/gin"
  "github.com/golang-jwt/jwt/v5"
)

type UserHandler struct {
  svc            *service.UserService
  codeSvc        *service.CodeService
  passwordRegexp *regexp.Regexp
  emailRegexp    *regexp.Regexp
  idRegexp       *regexp.Regexp
  phoneExp       *regexp.Regexp
}

const biz = "login"

var (
  ErrDuplicateUser = service.ErrDuplicateUser
  ErrBadRegExp     = errors.New("内部正则错误")
)

func NewUserHandler(svc *service.UserService, codeSvc *service.CodeService) *UserHandler {
  // 限制 8-50 长度的密码，防止加密算法不支持
  const passwordExp = `^(?=.*[a-z])(?=.*[A-Z])(?=.*[\d])[~!@#$%^&a-zA-Z\d]{8,50}$`
  const emailExp = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
  const idExp = `^[0-9]+$`
  const phoneExp = `^1\d{10}`
  return &UserHandler{
    passwordRegexp: regexp.MustCompile(passwordExp, 0),
    emailRegexp:    regexp.MustCompile(emailExp, 0),
    idRegexp:       regexp.MustCompile(idExp, 0),
    phoneExp:       regexp.MustCompile(phoneExp, 0),
    svc:            svc,
    codeSvc:        codeSvc,
  }
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
  group := server.Group("/user")
  group.POST("/register", u.Register)
  group.POST("/login", u.Login)
  group.POST("/logout", u.Logout)
  group.POST("/edit", u.Edit)
  group.GET("/profile", u.Profile)
  group.GET("/profileJwt", u.ProfileByJwt)
  group.POST("/sms/login/code", u.SendLoginSMSCode)
  group.POST("/sms/login/verify", u.VerifyLoginSMSCode)
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
  //session := sessions.Default(ctx)
  //session.Set("userId", user.Id)
  //session.Options(sessions.Options{
  //  MaxAge:   60,
  //  Path:     "/",
  //  HttpOnly: true,
  //})
  //session.Save()

  // JWT 版本
  // 1. 初始化jwt，带数据，不带入敏感数据（密码）
  // token := jwt.New(jwt.SigningMethodHS256)
  err = u.setJWTToken(ctx, user)
  if err != nil {
    ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "token生成错误"})
    return
  }

  // 登录成功返回用户数据
  ctx.JSON(http.StatusOK, gin.H{
    "msg":     "登录成功",
    "profile": parseUser(user),
  })
  return
}

func (u *UserHandler) setJWTToken(ctx *gin.Context, user domain.User) error {
  userClaim := UserClaim{
    Uid: user.Id,
    RegisteredClaims: jwt.RegisteredClaims{
      // jwt 过期时间
      // 60m 过期
      ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
    },
    UserAgent: ctx.Request.UserAgent(), // ua 用来判断是否异地登录
  }

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim)
  // 2. 加签，返回token字符
  tokenStr, err := token.SignedString(variable.JWTEncryptKey)
  if err != nil {
    return err
  }
  ctx.Header("x-jwt-token", tokenStr) // 返回前端存储，下次请求通过请求头带上来
  return nil
}

func (u *UserHandler) Logout(ctx *gin.Context) {
  // 退出登录 => 将登录失效 => 登录cookie清除
  sess := sessions.Default(ctx)
  // 指定和Login对应的path，才能删除cookie
  sess.Options(sessions.Options{
    MaxAge: -1,
    Path:   "/",
  })
  sess.Save()

  // 可以删除host的 cookie
  //ctx.SetCookie("wei_session", "", -1, "/", "localhost", false, false)
  ctx.JSON(http.StatusOK, gin.H{
    "msg": "退出成功",
  })
  return
}

func (u *UserHandler) Edit(ctx *gin.Context) {
  type Req struct {
    Birthday string `json:"birthday"`
  }
  var req Req
  if err := ctx.Bind(&req); err != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
    return
  }
  date, err := time.Parse(time.DateOnly, req.Birthday)
  if err != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": "请输入正确格式日期"})
    return
  }
  userId := sessions.Default(ctx).Get("userId").(int64)
  user, err := u.svc.UpdateProfile(ctx, domain.User{Birthday: date, Id: userId})
  if errors.Is(err, service.ErrOperateFail) {
    ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
    return
  }
  type User struct {
    Id       int64
    Name     string
    Email    string
    Birthday string
  }
  ctx.JSON(http.StatusOK, gin.H{
    "profile": parseUser(user),
  })
  return
}

// session 版本
func (u *UserHandler) Profile(ctx *gin.Context) {
  // 通过 id 查找用户
  type Req struct {
    Id int64 `form:"id"`
  }
  var req Req
  err := ctx.Bind(&req)
  if err != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
    return
  }

  id := strconv.FormatInt(req.Id, 10)
  ok, err := u.idRegexp.MatchString(id)
  if err != nil {
    ctx.JSON(http.StatusInternalServerError, ErrBadRegExp)
    return
  }
  if !ok {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
    return
  }

  // 这里 id 也可以通过 login 中间件设置的ctx.id获取
  user, err := u.svc.FindById(ctx, req.Id)
  if errors.Is(err, service.ErrRecordNoFound) {
    ctx.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
    return
  }

  ctx.JSON(http.StatusOK, gin.H{
    "profile": parseUser(user),
  })
}

// jwt 版本
func (u *UserHandler) ProfileByJwt(ctx *gin.Context) {
  userInfo, ok := ctx.Get("userInfo")
  if !ok {
    // 追踪没有存入的数据
    ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "系统错误"})
    return
  }
  uc, ok := userInfo.(UserClaim)
  if !ok {
    // 断言失败
    ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "系统错误"})
    return
  }
  user, err := u.svc.FindById(ctx, uc.Uid)
  if errors.Is(err, service.ErrRecordNoFound) {
    ctx.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
    return
  }

  ctx.JSON(http.StatusOK, gin.H{
    "profile": parseUser(user),
  })
}

func (u *UserHandler) LoginSMS(ctx *gin.Context) {
  type Req struct {
    Phone string `json:"phone"`
    Code  string `json:"code"`
  }
  var req Req
  if err := ctx.ShouldBind(&req); err != nil {
    ctx.JSON(http.StatusOK, res.Result{
      Code: 4,
      Msg:  "请输入正确的手机号/验证码",
    })
    return
  }
  // todo 手机号校验
  ok, err := u.codeSvc.Verify(ctx, biz, req.Code, req.Phone)
  if err != nil {
    ctx.JSON(http.StatusOK, res.Result{
      Code: 5,
      Msg:  "系统错误",
    })
    return
  }
  if !ok {
    ctx.JSON(http.StatusOK, res.Result{
      Code: 4,
      Msg:  "请输入正确的手机号/验证码",
    })
    return
  }

  // 校验成功
  // 有可能是新用户，通过手机号注册
  user, err := u.svc.FindOrCreate(ctx, req.Phone)
  if err != nil {
    ctx.JSON(http.StatusOK, res.Result{
      Code: 5,
      Msg:  "系统错误",
    })
    return
  }

  // 创建新的用户，设置 jwt 登录状态
  err = u.setJWTToken(ctx, user)
  if err != nil {
    ctx.JSON(http.StatusOK, res.Result{
      Code: 5,
      Msg:  "系统错误",
    })
    return
  }

  // 登录成功
  ctx.JSON(http.StatusOK, res.Result{
    Code: 0,
    Msg:  "验证码校验成功",
  })
}

func (u *UserHandler) SendLoginSMSCode(ctx *gin.Context) {
  const biz = "login"
  type Req struct {
    Phone string `json:"phone"`
  }
  var req Req
  if err := ctx.ShouldBindJSON(&req); err != nil {
    ctx.JSON(http.StatusBadRequest, res.Result{
      Code: 4,
      Msg:  "请输入正确的手机号/验证码",
    })
    return
  }
  ok, err := u.phoneExp.MatchString(req.Phone)
  if err != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{"error": "系统内部错误"})
    return
  }
  if !ok {
    ctx.JSON(http.StatusBadRequest, res.Result{
      Code: 4,
      Msg:  "请输入正确的手机号/验证码",
    })
    return
  }
  err = u.codeSvc.Send(ctx, biz, req.Phone)
  switch err {
  case nil:
    ctx.JSON(http.StatusOK, res.Result{
      Msg: "发送成功",
    })
  case service.ErrCodeSendTooMany:
    ctx.JSON(http.StatusOK, res.Result{
      Msg: "发送太频繁",
    })
  default:
    ctx.JSON(http.StatusOK, res.Result{
      Msg: "系统错误",
    })
    return
  }
  ctx.JSON(http.StatusOK, gin.H{
    "code": "Success",
  })
}

func (u *UserHandler) VerifyLoginSMSCode(ctx *gin.Context) {
  type Req struct {
    Phone string `json:"phone"`
    Code  string `json:"code"`
  }
  var req Req
  if err := ctx.ShouldBind(&req); err != nil {
    ctx.JSON(http.StatusBadRequest, res.Result{
      Code: 4,
      Msg:  "请输入正确的手机号/验证码",
    })
    return
  }
  valid, err := u.codeSvc.Verify(ctx, biz, req.Code, req.Phone)
  if err != nil {
    ctx.JSON(http.StatusOK, res.Result{
      Code: 5,
      Msg:  "系统错误",
    })
    return
  }
  if !valid {
    ctx.JSON(http.StatusBadRequest, res.Result{
      Code: 4,
      Msg:  "请输入正确的手机号/验证码",
    })
    return
  }
  // 验证码登录
  fmt.Println("sms login=====", req.Phone)
}

type User struct {
  Id       int64
  Name     string
  Email    string
  Birthday string
}

func parseUser(user domain.User) User {
  return User{
    Id:       user.Id,
    Name:     user.Name,
    Email:    user.Email,
    Birthday: user.Birthday.Format(time.DateOnly),
  }
}

// jwt 中 用户信息的 声明
type UserClaim struct {
  jwt.RegisteredClaims
  Uid       int64
  UserAgent string
}
