package repo

import (
  "context"
  "database/sql"
  "testing"
  "time"
  "weibook/internal/domain"
  "weibook/internal/repo/cache"
  cachemocks "weibook/internal/repo/cache/mock"
  "weibook/internal/repo/dao"
  daomocks "weibook/internal/repo/dao/mock"

  "github.com/stretchr/testify/assert"
  "go.uber.org/mock/gomock"
)

func TestCacheUserRepo_FindById(t *testing.T) {
  testCases := []struct {
    name string
    mock func(t *gomock.Controller) (dao.UserDAO, cache.UserCache)

    id int64

    wantUser domain.User
    wantErr  error
  }{
    {
      name: "缓存不中，数据库查询成功",
      mock: func(t *gomock.Controller) (dao.UserDAO, cache.UserCache) {
        uc := cachemocks.NewMockUserCache(t)
        uc.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.User{}, cache.ErrUserNotFound)

        ud := daomocks.NewMockUserDAO(t)
        ud.EXPECT().FindById(gomock.Any(), int64(123)).Return(dao.User{
          Id: 123,
          Email: sql.NullString{
            Valid:  true,
            String: "123123@qq.com",
          },
          Birthday: 0,
        }, nil)

        // 如果设置缓存不使用 goroutine则需要补充下面测试
        uc.EXPECT().Set(gomock.Any(), domain.User{
          Id:       123,
          Email:    "123123@qq.com",
          Birthday: time.UnixMilli(0),
        }).Return(nil)

        return ud, uc
      },

      id: 123,

      wantErr: nil,
      wantUser: domain.User{
        Id:       123,
        Email:    "123123@qq.com",
        Birthday: time.UnixMilli(0),
      },
    },
    {
      name: "缓存命中",
      mock: func(t *gomock.Controller) (dao.UserDAO, cache.UserCache) {
        uc := cachemocks.NewMockUserCache(t)
        uc.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.User{
          Id:       123,
          Email:    "123123@qq.com",
          Birthday: time.UnixMilli(0),
        }, nil)

        ud := daomocks.NewMockUserDAO(t)

        return ud, uc
      },

      id: 123,

      wantErr: nil,
      wantUser: domain.User{
        Id:       123,
        Email:    "123123@qq.com",
        Birthday: time.UnixMilli(0),
      },
    },
    {
      name: "查询失败",
      mock: func(t *gomock.Controller) (dao.UserDAO, cache.UserCache) {
        uc := cachemocks.NewMockUserCache(t)
        uc.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.User{}, cache.ErrUserNotFound)

        ud := daomocks.NewMockUserDAO(t)
        ud.EXPECT().FindById(gomock.Any(), int64(123)).Return(dao.User{}, dao.ErrUserNotFound)

        return ud, uc
      },

      id: 123,

      wantErr:  ErrUserNotFound,
      wantUser: domain.User{},
    },
  }

  for _, tt := range testCases {
    t.Run(tt.name, func(t *testing.T) {
      c := gomock.NewController(t)
      defer c.Finish()
      ud, uc := tt.mock(c)
      repo := NewUserRepo(ud, uc)
      user, err := repo.FindById(context.Background(), tt.id)
      assert.Equal(t, err, tt.wantErr)
      assert.Equal(t, user, tt.wantUser)

      // 使用 goroutine
      time.Sleep(time.Second)
    })
  }
}
