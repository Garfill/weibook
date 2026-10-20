//go:build wireinject

package main

import (
  "weibook/internal/repo"
  "weibook/internal/repo/cache"
  "weibook/internal/repo/dao"
  "weibook/internal/service"
  web "weibook/internal/web/user"
  "weibook/ioc"

  "github.com/gin-gonic/gin"
  "github.com/goforj/wire"
)

var AppSet = wire.NewSet(
  // db
  ioc.InitDB,
  ioc.InitRedis,
  // dao
  dao.NewUserDAO,
  cache.NewUserCache,
  cache.NewCodeCache,
  // repo
  repo.NewUserRepo,
  repo.NewCodeRepo,
  // service
  service.NewUserService,
  service.NewCodeService,
  ioc.InitSmsService,
  // handler
  web.NewUserHandler,
  // middleware
  ioc.InitMiddlewares,
  // gin
  ioc.InitGin,
)

func InitializeApp() *gin.Engine {
  wire.Build(AppSet)
  return nil
}
