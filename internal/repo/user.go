package repo

import (
  "context"
  "database/sql"
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

type UserRepo interface {
  CreateUser(ctx context.Context, u domain.User) error
  FindByEmail(ctx context.Context, email string) (domain.User, error)
  FindByPhone(ctx context.Context, phone string) (domain.User, error)
  FindById(ctx context.Context, id int64) (domain.User, error)
  Update(ctx context.Context, profile domain.User) (domain.User, error)
  domainToEntity(u domain.User) dao.User
  entityToDomain(ud dao.User) domain.User
}

type CacheUserRepo struct {
  dao   dao.UserDAO
  cache cache.UserCache
}

func NewUserRepo(dao dao.UserDAO, ca cache.UserCache) UserRepo {
  return &CacheUserRepo{
    dao:   dao,
    cache: ca,
  }
}

func (repo *CacheUserRepo) CreateUser(ctx context.Context, u domain.User) error {
  return repo.dao.Insert(ctx, repo.domainToEntity(u))
}

func (repo *CacheUserRepo) FindByEmail(ctx context.Context, email string) (domain.User, error) {
  u, err := repo.dao.FindByEmail(ctx, email)
  if errors.Is(err, ErrUserNotFound) {
    return domain.User{}, ErrUserNotFound
  }
  return repo.entityToDomain(u), nil
}

func (repo *CacheUserRepo) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
  u, err := repo.dao.FindByPhone(ctx, phone)
  if errors.Is(err, ErrUserNotFound) {
    return domain.User{}, ErrUserNotFound
  }
  return repo.entityToDomain(u), nil
}

func (repo *CacheUserRepo) FindById(ctx context.Context, id int64) (domain.User, error) {
  // 1. 先从缓存取数据
  user, err := repo.cache.Get(ctx, id)
  // 缓存有/没有/出错，不知道有没有数据返回
  if err == nil {
    // 有数据
    return user, err
  }
  if err == cache.ErrUserNotFound {
    // 缓存没数据
    // 2. 从数据库查找
    ue, err := repo.dao.FindById(ctx, id)
    if errors.Is(err, ErrUserNotFound) {
      return domain.User{}, ErrUserNotFound
    }
    if err != nil {
      return domain.User{}, ErrOperateFail
    }
    user = repo.entityToDomain(ue)
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

func (repo *CacheUserRepo) Update(ctx context.Context, profile domain.User) (domain.User, error) {
  u, err := repo.dao.Update(ctx, profile)
  if err != nil {
    return domain.User{}, err
  }
  return repo.entityToDomain(u), nil
}

// dao和 domain 的 email/phone 类型不一致，转换函数
func (repo *CacheUserRepo) domainToEntity(u domain.User) dao.User {
  return dao.User{
    Id:   u.Id,
    Name: u.Name,
    Email: sql.NullString{
      String: u.Email,
      Valid:  u.Email != "",
    },
    Phone: sql.NullString{
      String: u.Phone,
      Valid:  u.Phone != "",
    },
    Password: u.Password,
    Birthday: u.Birthday.UnixMilli(),
  }
}

func (repo *CacheUserRepo) entityToDomain(ud dao.User) domain.User {
  return domain.User{
    Id:       ud.Id,
    Name:     ud.Name,
    Email:    ud.Email.String,
    Phone:    ud.Phone.String,
    Password: ud.Password,
    Birthday: time.UnixMilli(ud.Birthday).In(time.Local),
  }
}
