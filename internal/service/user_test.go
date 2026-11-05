package service

import (
  "context"
  "testing"
  "weibook/internal/domain"
  "weibook/internal/repo"
  repomocks "weibook/internal/repo/mock"

  "github.com/stretchr/testify/assert"
  "go.uber.org/mock/gomock"
)

func TestUserSvc_Login(t *testing.T) {
  testCases := []struct {
    name     string
    mock     func(c *gomock.Controller) repo.UserRepo
    email    string
    password string
    wantErr  error
    wantUser domain.User
  }{
    {
      name: "登录成功",
      mock: func(c *gomock.Controller) repo.UserRepo {
        userRepo := repomocks.NewMockUserRepo(c)
        userRepo.EXPECT().
          FindByEmail(gomock.Any(), "12345@qq.com").
          Return(domain.User{
            Email: "12345@qq.com",
            // 你在这边拿到的密码，就应该是一个正确的密码
            // 加密后的正确的密码
            Password: "$2a$10$.l0JHmM7a2PdJ.A9gsmVyerEDlp1WhxsglC34S4UJH4TuHhWY7Tfq",
            Phone:    "15212345678",
          }, nil)
        return userRepo
      },

      email:    "12345@qq.com",
      password: "123456#hello",

      wantUser: domain.User{
        Email:    "12345@qq.com",
        Password: "$2a$10$.l0JHmM7a2PdJ.A9gsmVyerEDlp1WhxsglC34S4UJH4TuHhWY7Tfq",
        Phone:    "15212345678",
      },
      wantErr: nil,
    },
    {
      name: "用户不存在",
      mock: func(c *gomock.Controller) repo.UserRepo {
        userRepo := repomocks.NewMockUserRepo(c)
        userRepo.EXPECT().
          FindByEmail(gomock.Any(), "12345@qq.com").
          Return(domain.User{}, repo.ErrUserNotFound) // repo 层模拟返回的repo错误
        return userRepo
      },

      email:    "12345@qq.com",
      password: "123456#hello",

      wantUser: domain.User{},
      wantErr:  ErrInvalidUserOrPwd, // service 层接收到 repo 层错误，转化返回的service层结果错误
    },
    {
      name: "密码错误",
      mock: func(c *gomock.Controller) repo.UserRepo {
        userRepo := repomocks.NewMockUserRepo(c)
        userRepo.EXPECT().
          FindByEmail(gomock.Any(), "12345@qq.com").
          Return(domain.User{
            Email: "12345@qq.com",
            // 你在这边拿到的密码，就应该是一个正确的密码
            // 加密后的正确的密码
            Password: "$2a$10$.l0JHmM7a2PdJ.A9gsmVyerEDlp1WhxsglC34S4UJH4TuHhWY7Tfq",
            Phone:    "15212345678",
          }, nil)
        return userRepo
      },

      email:    "12345@qq.com",
      password: "1fdafado",

      wantUser: domain.User{},
      wantErr:  ErrInvalidUserOrPwd,
    },
  }

  for _, tc := range testCases {
    t.Run(tc.name, func(t *testing.T) {
      ctrl := gomock.NewController(t)
      defer ctrl.Finish()

      svc := NewUserService(tc.mock(ctrl))
      user, err := svc.Login(context.Background(), tc.email, tc.password)
      assert.Equal(t, tc.wantErr, err)
      assert.Equal(t, tc.wantUser, user)
    })
  }
}
