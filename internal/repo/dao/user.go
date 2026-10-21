package dao

import (
  "context"
  "database/sql"
  "errors"
  "weibook/internal/domain"

  "github.com/go-sql-driver/mysql"
  "gorm.io/gorm"
  "gorm.io/plugin/soft_delete"
)

var (
  ErrDuplicateUser = errors.New("手机/邮箱冲突")
  ErrUserNotFound  = gorm.ErrRecordNotFound
)

type UserDAO interface {
  Insert(ctx context.Context, u User) error
  FindByEmail(ctx context.Context, email string) (User, error)
  FindByPhone(ctx context.Context, phone string) (User, error)
  FindById(ctx context.Context, id int64) (User, error)
  Update(ctx context.Context, user domain.User) (User, error)
}

type GORMUserDAO struct {
  db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDAO {
  return &GORMUserDAO{
    db: db,
  }
}

func (dao *GORMUserDAO) Insert(ctx context.Context, u User) error {
  // WithContext 用来保持链路
  err := dao.db.WithContext(ctx).Create(&u).Error
  // 类型断言是mysql错误
  //if mysqlErr, ok := err.(*mysql.MySQLError); ok {
  //  const uniqueConflictsErrorNo = 1062 // 唯一索引冲突
  //  if mysqlErr.Number == uniqueConflictsErrorNo {
  //    return ErrDuplicateUser
  //  }
  //}

  if err != nil {
    var mysqlErr *mysql.MySQLError
    const uniqueConflictsErrNo = 1062
    if errors.As(err, &mysqlErr) && mysqlErr.Number == uniqueConflictsErrNo {
      // 唯一索引冲突
      return ErrDuplicateUser
    }
    return err
  }
  return nil
}

func (dao *GORMUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
  var u User
  err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
  return u, err
}

func (dao *GORMUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
  var u User
  err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
  return u, err
}

func (dao *GORMUserDAO) FindById(ctx context.Context, id int64) (User, error) {
  var u User
  err := dao.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
  return u, err
}

func (dao *GORMUserDAO) Update(ctx context.Context, user domain.User) (User, error) {
  u, err := dao.FindById(ctx, user.Id)
  if err != nil {
    return User{}, err
  }
  u.Birthday = user.Birthday.UnixMilli()
  dao.db.Save(&u)
  return u, nil
}

// 对标数据库内部的字段
// 别名 entity, model, PO(peristent object)
type User struct {
  Id       int64          `gorm:"primaryKey,autoIncrement"`
  Name     string         `gorm:"size:100;not null"`
  Password string         `gorm:"size:100;not null"`
  Email    sql.NullString `gorm:"index:,unique;size:100"` // 唯一索引可以有多个空值，但不能有多个空的字符""
  Phone    sql.NullString `gorm:"index:,unique;size:20"`
  Birthday int64

  // 时间存 时间戳不受时区影响
  CreatedAt int64                 `gorm:"autoCreateTime:milli"`
  UpdatedAt int64                 `gorm:"autoUpdateTime:milli"`
  DeletedAt soft_delete.DeletedAt `gorm:"softDelete:milli"`
}
