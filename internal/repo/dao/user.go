package dao

import (
  "context"
  "errors"
  "github.com/go-sql-driver/mysql"
  "gorm.io/gorm"
)

var DuplicateUserEmailErr = errors.New("邮箱冲突")

type UserDAO struct {
  db *gorm.DB
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
  // WithContext 用来保持链路
  err := dao.db.WithContext(ctx).Create(&u).Error
  // 类型断言是mysql错误
  if mysqlErr, ok := err.(*mysql.MySQLError); ok {
    const uniqueConflictsErrorNo = 1062 // 唯一索引冲突
    if mysqlErr.Number == uniqueConflictsErrorNo {
      return DuplicateUserEmailErr
    }
  }
  return err
}

func NewUserDAO(db *gorm.DB) *UserDAO {
  return &UserDAO{
    db: db,
  }
}

// 对标数据库内部的字段
// 别名 entity, model, PO(peristent object)
type User struct {
  Id       uint64 `gorm:"primaryKey,autoIncrement"`
  Name     string `gorm:"size:100;not null"`
  Password string `gorm:"size:100;not null"`
  Email    string `gorm:"index:,unique;size:100"`

  // 时间存 时间戳不受时区影响
  CreaeteAt int64 `gorm:"autoCreateTime:milli"`
  UpdateAt  int64 `gorm:"autoUpdateTime:milli"`
  DeletedAt int64
}
