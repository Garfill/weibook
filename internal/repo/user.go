package repo

import (
  "context"
  "errors"
  "time"
  "weibook/internal/domain"
  "weibook/internal/repo/cache"
  "weibook/internal/repo/dao"
)

var (
  ErrDuplicateUser = dao.ErrDuplicateUser
  ErrUserNotFound  = dao.ErrUserNotFound
  ErrOperateFail   = errors.New("操作失败")
)

type UserRepo struct {
  dao   *dao.UserDAO
  cache *cache.UserCache
}

func NewUserRepo(dao *dao.UserDAO, ca *cache.UserCache) *UserRepo {
  return &UserRepo{
    dao:   dao,
    cache: ca,
  }
}

func (repo *UserRepo) CreateUser(ctx context.Context, u domain.User) error {
  return repo.dao.Insert(ctx, dao.User{
    Name:     u.Name,
    Password: u.Password,
    Email:    u.Email,
  })
}

func (repo *UserRepo) FindByEmail(ctx context.Context, email string) (domain.User, error) {
  u, err := repo.dao.FindByEmail(ctx, email)
  if err != nil {
    return domain.User{}, err
  }
  return domain.User{
    Id:       u.Id,
    Email:    u.Email,
    Name:     u.Name,
    Password: u.Password,
    Birthday: time.UnixMilli(u.Birthday),
  }, nil
}

func (repo *UserRepo) FindById(ctx context.Context, id int64) (domain.User, error) {
  // 1. 先从缓存取数据
  user, err := repo.cache.Get(ctx, id)
  // 缓存有/没有/出错，不知道有没有数据返回
  if err == nil {
    // 有数据
    return user, err
  }
  if err == cache.ErrUserNotFound {
    //  缓存没数据
    // 2. 从数据库查找
    ue, err := repo.dao.FindById(ctx, id)
    if errors.Is(err, ErrUserNotFound) {
      return domain.User{}, ErrUserNotFound
    }
    if err != nil {
      return domain.User{}, ErrOperateFail
    }
    user = domain.User{
      Id:       ue.Id,
      Email:    ue.Email,
      Name:     ue.Name,
      Password: ue.Password,
      Birthday: time.UnixMilli(ue.Birthday),
    }
    go func() {
      //  设置缓存数据
      err = repo.cache.Set(ctx, user)
      if err != nil {
        // 可以日志监控
        // 缓存设置失败
      }
    }()
    return user, nil
  }

  // 其他错误 io.EOF / redis 崩溃 / 可能偶发性错误
  // 如果执行数据库加载，可能缓存击穿，搭配数据库限流 / 布尔过滤器

  // 这里考虑redis服务不可用
  return domain.User{}, nil
}

func (repo *UserRepo) Update(ctx context.Context, profile domain.User) (domain.User, error) {
  u, err := repo.dao.Update(ctx, profile)
  if err != nil {
    return domain.User{}, err
  }
  return domain.User{
    Id:       u.Id,
    Name:     u.Name,
    Email:    u.Email,
    Birthday: time.UnixMilli(u.Birthday),
  }, nil
}
