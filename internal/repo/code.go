package repo

import (
  "context"
  "weibook/internal/repo/cache"
)

var (
  ErrSendTooMany   = cache.ErrCodeSendTooMany
  ErrVerifyTooMany = cache.ErrCodeVerifyTooMany
)

type CodeRepo struct {
  cache *cache.CodeCache
}

func NewCodeRepo(c *cache.CodeCache) *CodeRepo {
  return &CodeRepo{
    cache: c,
  }
}

func (repo *CodeRepo) Store(ctx context.Context, biz, phone, code string) error {
  return repo.cache.Set(ctx, biz, phone, code)
}

func (repo *CodeRepo) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
  return repo.cache.Verify(ctx, biz, phone, code)
}
