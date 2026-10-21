package repo

import (
  "context"
  "weibook/internal/repo/cache"
)

var (
  ErrSendTooMany   = cache.ErrCodeSendTooMany
  ErrVerifyTooMany = cache.ErrCodeVerifyTooMany
)

type CodeRepo interface {
  Store(ctx context.Context, biz, phone, code string) error
  Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

type CacheCodeRepo struct {
  cache cache.CodeCache
}

func NewCodeRepo(c cache.CodeCache) CodeRepo {
  return &CacheCodeRepo{
    cache: c,
  }
}

func (repo *CacheCodeRepo) Store(ctx context.Context, biz, phone, code string) error {
  return repo.cache.Set(ctx, biz, phone, code)
}

func (repo *CacheCodeRepo) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
  return repo.cache.Verify(ctx, biz, phone, code)
}
