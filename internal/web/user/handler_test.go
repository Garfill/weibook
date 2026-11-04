package web

import (
  "bytes"
  "context"
  "encoding/json"
  "errors"
  "net/http"
  "net/http/httptest"
  "testing"
  "weibook/internal/domain"
  "weibook/internal/service"
  svcmocks "weibook/internal/service/mock"

  "github.com/gin-gonic/gin"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/require"
  "go.uber.org/mock/gomock"
)

func TestUserHandler_Register(t *testing.T) {
  testCases := []struct {
    name     string
    mock     func(c *gomock.Controller) service.UserService
    body     string
    wantCode int
    errorMsg string
  }{
    {
      name: "注册成功",
      mock: func(c *gomock.Controller) service.UserService {
        svc := svcmocks.NewMockUserService(c)
        svc.EXPECT().SignUp(gomock.Any(), domain.User{
          Name:     "guanjunhui",
          Password: "guan@1111",
          Email:    "12345@qq.com",
        }).
          Return(nil) // 注册成功返回nil
        return svc
      },
      body: `
{
  "name": "guanjunhui",
  "password": "guan@1111",
  "email": "12345@qq.com"
}
`,
      wantCode: http.StatusOK,
    },
    {
      name: "邮箱格式错误",
      mock: func(c *gomock.Controller) service.UserService {
        svc := svcmocks.NewMockUserService(c)
        // 参数错误情况下不需要调用signup
        return svc
      },
      body: `
    {
     "name": "guanjunhui",
     "password": "guan@1111",
      "email": "123"
    }
    `,
      wantCode: http.StatusBadRequest,
      errorMsg: "请输入正确格式的邮箱",
    },
    {
      name: "密码格式错误",
      mock: func(c *gomock.Controller) service.UserService {
        svc := svcmocks.NewMockUserService(c)
        // 参数错误情况下不需要调用signup
        return svc
      },
      body: `
    {
     "name": "guanjunhui",
     "password": "111",
      "email": "12345@qq.com"
    }
    `,
      wantCode: http.StatusBadRequest,
      errorMsg: "密码格式错误",
    },
    {
      name: "重复邮箱",
      mock: func(c *gomock.Controller) service.UserService {
        svc := svcmocks.NewMockUserService(c)
        svc.EXPECT().SignUp(gomock.Any(), domain.User{
          Name:     "guanjunhui",
          Password: "guan@1111",
          Email:    "12345@qq.com",
        }).
          Return(service.ErrDuplicateUser) // 注册成功返回nil
        return svc
      },
      body: `
{
  "name": "guanjunhui",
  "password": "guan@1111",
  "email": "12345@qq.com"
}
`,
      wantCode: http.StatusBadRequest,
      errorMsg: "手机/邮箱冲突",
    },
    {
      name: "系统错误",
      mock: func(c *gomock.Controller) service.UserService {
        svc := svcmocks.NewMockUserService(c)
        svc.EXPECT().SignUp(gomock.Any(), domain.User{
          Name:     "guanjunhui",
          Password: "guan@1111",
          Email:    "12345@qq.com",
        }).
          Return(errors.New("NotOk")) // 注册成功返回nil
        return svc
      },
      body: `
{
  "name": "guanjunhui",
  "password": "guan@1111",
  "email": "12345@qq.com"
}
`,
      wantCode: http.StatusInternalServerError,
      errorMsg: "系统错误",
    },
  }

  // 单元测试，不连接数据库等第三方，使用uber-mock实现

  for _, tc := range testCases {
    t.Run(tc.name, func(t *testing.T) {
      // 运行所有测试用例
      // go的json串里最后一项不加逗号！！！
      req, err := http.NewRequest(http.MethodPost, "/user/register", bytes.NewBuffer([]byte(tc.body)))
      req.Header.Set("Content-Type", "application/json")
      // 自己构造的请求不会有error
      require.NoError(t, err)

      // 使用 recorder 记录结果
      resp := httptest.NewRecorder()

      ctrl := gomock.NewController(t)
      defer ctrl.Finish()
      handler := NewUserHandler(tc.mock(ctrl), nil)
      server := gin.Default()
      handler.RegisterRoutes(server)
      server.ServeHTTP(resp, req) // 调用server处理请求入口
      assert.Equal(t, tc.wantCode, resp.Code)

      if resp.Code != http.StatusOK {
        // 应对实际业务中，返回 StatusCode 均为200，业务码code在body中的情况
        type Res struct {
          Error string `json:"error"`
        }
        var result Res
        json.Unmarshal(resp.Body.Bytes(), &result)
        assert.Equal(t, tc.errorMsg, result.Error)
      }

    })
  }
}

func TestMock(t *testing.T) {
  ctrl := gomock.NewController(t)
  defer ctrl.Finish()

  usersvc := svcmocks.NewMockUserService(ctrl)
  // 预期输入
  usersvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).
    Return(errors.New("mock error"))

  // 测试调用
  err := usersvc.SignUp(context.Background(), domain.User{
    Email: "12345678@qq.com",
  })
  t.Log(err)
}
