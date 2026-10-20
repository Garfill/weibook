package main

import (
  "weibook/internal/config"

  "github.com/gin-contrib/sessions"
  sessionRedis "github.com/gin-contrib/sessions/redis"
  "github.com/gin-gonic/gin"
  "github.com/redis/go-redis/v9"
  "github.com/ulule/limiter/v3"
  limiterGin "github.com/ulule/limiter/v3/drivers/middleware/gin"
  limiterRedis "github.com/ulule/limiter/v3/drivers/store/redis"
)

func main() {
  //初始化gin
  server := InitializeApp()
  // session
  initSession(server)
  initRateLimit(server)

  //server.GET("/hello", func(context *gin.Context) {
  //  context.JSON(200, gin.H{
  //    "message": "world",
  //  })
  //})
  server.Run(":8080")
  println("main ==========")
}

func initSession(server *gin.Engine) {
  // V1 版本，数据设置在cookie 不安全
  //store := cookie.NewStore([]byte("secret"))

  // 信息安全：
  // 1. 身份认证key和数据加密key
  // 2. 数据操作权限
  // 注意这里的key要使用随机生成的，除了保存到redis还会保存到cookie，不过cookie里没有数据只有sid

  // redis实现 V1 多实例部署
  store, _ := sessionRedis.NewStore(
    10,
    "tcp", config.Config.Redis.Addr, "", "",
    []byte("OEnEc62tqMFBOYRQEWQKmFWBvcpViJHV"), []byte("5H5v7Qqhct6EQBZ0DfsibYwi1J2l52xh"))
  sessionRedis.SetKeyPrefix(store, "wei_session_") // redis 内 key 前缀

  // memStore实现 V2 单机部署，基于内存的实现
  //store := memstore.NewStore([]byte("OEnEc62tqMFBOYRQEWQKmFWBvcpViJHV"))

  // 注册session中间件
  server.Use(sessions.Sessions("weibook_session", store))

  // 自定义中间件
  //server.Use(middleware.NewLoginMiddleBuilder().Build())
  //server.Use(middleware.NewLoginJWTMiddleBuilder().Build())
}

func initRateLimit(server *gin.Engine) {
  // 限流
  rate, err := limiter.NewRateFromFormatted("100-S")
  if err != nil {
    panic("Failed to initialize rate limiter")
  }
  client := redis.NewClient(&redis.Options{
    Addr: config.Config.Redis.Addr,
  })
  store, err := limiterRedis.NewStoreWithOptions(client, limiter.StoreOptions{
    Prefix: "limiter_prefix_",
  })
  if err != nil {
    panic("Failed to initialize limit store")
  }
  mw := limiterGin.NewMiddleware(limiter.New(store, rate))
  server.Use(mw)
}
