package dao

import (
  "context"
  "errors"

  "github.com/go-sql-driver/mysql"
  "gorm.io/gorm"
)

var (
  ErrDuplicateUser = errors.New("邮箱冲突")
  ErrUserNotFound  = gorm.ErrRecordNotFound
)

type UserDAO struct {
  db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
  return &UserDAO{
    db: db,
  }
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
  // WithContext 用来保持链路
  err := dao.db.WithContext(ctx).Create(&u).Error
  // 类型断言是mysql错误
  if mysqlErr, ok := err.(*mysql.MySQLError); ok {
    const uniqueConflictsErrorNo = 1062 // 唯一索引冲突
    if mysqlErr.Number == uniqueConflictsErrorNo {
      return ErrDuplicateUser
    }
  }
  return err
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
  var u User
  err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
  return u, err
}

func (dao *UserDAO) FindById(ctx context.Context, id string) (User, error) {
  var u User
  //strId, _ := strconv.ParseInt(id, 10, 64)
  err := dao.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
  return u, err
}

// 对标数据库内部的字段
// 别名 entity, model, PO(peristent object)
type User struct {
  Id       int64  `gorm:"primaryKey,autoIncrement"`
  Name     string `gorm:"size:100;not null"`
  Password string `gorm:"size:100;not null"`
  Email    string `gorm:"index:,unique;size:100"`

  // 时间存 时间戳不受时区影响
  CreatedAt int64 `gorm:"autoCreateTime:milli"`
  UpdatedAt int64 `gorm:"autoUpdateTime:milli"`
  DeletedAt int64
}
