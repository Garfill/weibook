package dao

import (
  "context"
  "gorm.io/gorm"
)

type UserDAO struct {
  db *gorm.DB
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
  // WithContext 用来保持链路
  return dao.db.WithContext(ctx).Create(&u).Error
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
  Name     string
  Password string

  // 时间存 时间戳不受时区影响
  CreaeteAt int64 `gorm:"autoCreateTime:milli"`
  UpdateAt  int64 `gorm:"autoUpdateTime:milli"`
  Deleted   gorm.DeletedAt
}
