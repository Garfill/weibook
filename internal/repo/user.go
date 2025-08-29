package repo

import (
  "context"
  "errors"
  "weibook/internal/domain"
  "weibook/internal/repo/dao"
)

var (
  ErrDuplicateUser = dao.ErrDuplicateUser
  ErrUserNotFound  = dao.ErrUserNotFound
)

type UserRepo struct {
  dao *dao.UserDAO
}

func NewUserRepo(dao *dao.UserDAO) *UserRepo {
  return &UserRepo{
    dao: dao,
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
  user, err := repo.dao.FindByEmail(ctx, email)
  if err != nil {
    return domain.User{}, err
  }
  return domain.User{
    Id:       user.Id,
    Email:    user.Email,
    Name:     user.Name,
    Password: user.Password,
  }, nil
}

func (repo *UserRepo) Update(ctx context.Context, profile domain.User) (domain.User, error) {
  u, _ := repo.dao.Update(ctx, profile)
  return domain.User{
    Id:       u.Id,
    Name:     u.Name,
    Birthday: u.Birthday,
    Email:    u.Email,
  }, nil
}

func (repo *UserRepo) FindById(ctx context.Context, id string) (domain.User, error) {
  user, err := repo.dao.FindById(ctx, id)
  if errors.Is(err, ErrUserNotFound) {
    return domain.User{}, ErrUserNotFound
  }
  if err != nil {
    return domain.User{}, err
  }
  return domain.User{
    Id:       user.Id,
    Email:    user.Email,
    Name:     user.Name,
    Password: user.Password,
  }, nil
}
