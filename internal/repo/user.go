package repo

import (
  "context"
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
    Password: user.Password,
  }, nil
}
