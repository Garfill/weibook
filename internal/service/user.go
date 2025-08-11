package service

import (
  "context"
  "errors"
  "golang.org/x/crypto/bcrypt"
  "weibook/internal/domain"
  "weibook/internal/repo"
)

type UserService struct {
  repo *repo.UserRepo
}

var (
  ErrDuplicateUser    = repo.ErrDuplicateUser
  ErrInvalidUserOrPwd = errors.New("帐号或者密码错误")
)

func NewUserService(repo *repo.UserRepo) *UserService {
  return &UserService{repo: repo}
}

func (svc *UserService) SignUp(ctx context.Context, user domain.User) error {
  // 1. context 是为了保持连续性和链路控制
  // 2. 使用实体 user 而不使用指针，可以不写 user == nil 判空（但是会触发复制 user）
  // 3. 实体 user 有利于分配到栈，不会内存逃逸
  // 4. 良好习惯 函数返回 error
  // 5. 这一层 考虑 数据加密 同时 执行保存操作
  hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
  if err != nil {
    return err
  }
  user.Password = string(hash)
  return svc.repo.CreateUser(ctx, user)
}

func (svc *UserService) Login(ctx context.Context, email string, password string) (domain.User, error) {
  // 查找用户
  u, err := svc.repo.FindByEmail(ctx, email)
  if errors.Is(err, repo.ErrUserNotFound) {
    return domain.User{}, ErrInvalidUserOrPwd
  }
  if err != nil {
    // 系统内部错误
    return domain.User{}, err
  }
  // 对比密码是否一致
  err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
  if err != nil {
    return domain.User{}, ErrInvalidUserOrPwd
  }
  // 将查找用户失败或者密码错误的error转化为同一个error返回
  return u, nil
}

// service 层级不定义 User，而是在 domain 层级 定义
// 参考面对对象编程，User 属于一个对象（这里叫领域对象）
