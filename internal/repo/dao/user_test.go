package dao

import (
  "context"
  "database/sql"
  "testing"

  "github.com/DATA-DOG/go-sqlmock"
  "github.com/go-sql-driver/mysql"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/require"
  gormMysql "gorm.io/driver/mysql"
  "gorm.io/gorm"
)

func TestGORMUserDAO_Insert(t *testing.T) {
  tcs := []struct {
    name string
    user User

    mock func(t *testing.T) *sql.DB

    wantErr error
  }{
    {
      name: "success",
      mock: func(t *testing.T) *sql.DB {
        mockDB, mock, err := sqlmock.New()
        require.NoError(t, err)
        mockRes := sqlmock.NewResult(1, 1)
        // 预期是正则表达式执行
        // 增删改用 ExpectExec, 查询用 ExpectQuery
        mock.ExpectExec("INSERT INTO `users` .*").
          WillReturnResult(mockRes)
        return mockDB
      },
      user: User{
        Id:       1,
        Name:     "1111111",
        Password: "123123",
      },
      wantErr: nil,
    },
    {
      name: "邮箱冲突",
      mock: func(t *testing.T) *sql.DB {
        mockDB, mock, err := sqlmock.New()
        require.NoError(t, err)
        // 预期是正则表达式执行
        // 增删改用 ExpectExec, 查询用 ExpectQuery
        mock.ExpectExec("INSERT INTO `users` .*").
          WillReturnError(&mysql.MySQLError{Number: 1062})
        return mockDB
      },
      user: User{
        Id:       1,
        Name:     "1111111",
        Password: "123123",
      },
      wantErr: ErrDuplicateUser,
    },
  }
  for _, tt := range tcs {
    t.Run(tt.name, func(t *testing.T) {
      mockDB := tt.mock(t)
      db, err := gorm.Open(gormMysql.New(gormMysql.Config{
        Conn: mockDB,
        // 跳过默认 select version
        SkipInitializeWithVersion: true,
      }), &gorm.Config{
        // mock 不用自动ping
        DisableAutomaticPing: true,
        // mockdb 不开启默认事务
        SkipDefaultTransaction: true,
      })
      d := NewUserDAO(db)

      err = d.Insert(context.Background(), tt.user)
      assert.Equal(t, tt.wantErr, err)
    })
  }
}
