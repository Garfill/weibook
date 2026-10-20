package service

import (
  "context"
  "errors"
  "weibook/internal/domain"
  "weibook/internal/repo"

  "golang.org/x/crypto/bcrypt"
)

type UserService struct {
  repo *repo.UserRepo
}

var (
  ErrDuplicateUser    = repo.ErrDuplicateUser
  ErrInvalidUserOrPwd = errors.New("帐号或者密码错误")
  ErrRecordNoFound    = errors.New("没找到数据")
  ErrOperateFail      = repo.ErrOperateFail
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

func (svc *UserService) FindById(ctx context.Context, id int64) (domain.User, error) {
  user, err := svc.repo.FindById(ctx, id)
  if errors.Is(err, repo.ErrUserNotFound) {
    return domain.User{}, ErrRecordNoFound
  }
  return user, nil
}

// 1. 查询手机
// 2. 否则注册用户，成功则返回
// 3. 否则（遇到并发问题），查询返回
func (svc *UserService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
  // 手机号没找到用户
  u, err := svc.repo.FindByPhone(ctx, phone)
  if err != repo.ErrUserNotFound {
    // err = nil
    // err = 其他错误
    return u, err
  }
  // 没对应手机号用户
  u = domain.User{
    Phone: phone,
  }
  err = svc.repo.CreateUser(ctx, u)
  if err != nil && err != repo.ErrDuplicateUser {
    // 创建失败，但是不是手机冲突，属于系统错误
    return domain.User{}, err
  }
  // 创建后重新查找库里面的数据
  // tips 有可能遇到主从延迟问题，如果遇到，createUser 返回新生成的id（int64）
  return svc.repo.FindByPhone(ctx, phone)
}

func (svc *UserService) UpdateProfile(ctx context.Context, profile domain.User) (domain.User, error) {
  return svc.repo.Update(ctx, profile)
}

// service 层级不定义 User，而是在 domain 层级 定义
// 参考面对对象编程，User 属于一个对象（这里叫领域对象）
