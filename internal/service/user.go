package service

import (
  "context"
  "weibook/internal/domain"
  "weibook/internal/repo"
)

type UserService struct {
  repo *repo.UserRepo
}

func NewUserService(repo *repo.UserRepo) *UserService {
  return &UserService{repo: repo}
}

func (svc *UserService) SignUp(ctx context.Context, user domain.User) error {
  // 1. context 是为了保持连续性和链路控制
  // 2. 使用实体 user 而不使用指针，可以不写 user == nil 判空（但是会触发复制 user）
  // 3. 实体 user 有利于分配到栈，不会内存逃逸
  // 4. 良好习惯 函数返回 error
  // 5. 这一层 考虑 数据加密 同时 执行保存操作
  return svc.repo.CreateUser(ctx, user)

}

// service 层级不定义 User，而是在 domain 层级 定义
// 参考面对对象编程，User 属于一个对象（这里叫领域对象）
