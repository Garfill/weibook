package repo

import (
  "context"
  "weibook/internal/domain"
  "weibook/internal/repo/dao"
)

var DuplicateUserEmailErr = dao.DuplicateUserEmailErr

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
