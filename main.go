package main

import (
  "time"
  "weibook/internal/config"
  "weibook/internal/repo"
  "weibook/internal/repo/dao"
  "weibook/internal/service"
  "weibook/internal/www/middleware"
  "weibook/internal/www/user"

  "github.com/gin-contrib/cors"
  "github.com/gin-contrib/sessions"
  sessionRedis "github.com/gin-contrib/sessions/redis"
  "github.com/gin-gonic/gin"
  "github.com/redis/go-redis/v9"
  "github.com/ulule/limiter/v3"
  limiterGin "github.com/ulule/limiter/v3/drivers/middleware/gin"
  limiterRedis "github.com/ulule/limiter/v3/drivers/store/redis"
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
)

func main() {
  // 数据库初始化
  db := initDB()
  // user 初始化
  userHandler := initUser(db)
  //初始化gin
  server := initServer()
  // session
  initSession(server)
  // 限流
  initRateLimit(server)

  // 注册handler
  userHandler.RegisterRoutes(server)

  //server.GET("/hello", func(context *gin.Context) {
  //  context.JSON(200, gin.H{
  //    "message": "world",
  //  })
  //})
  server.Run(":8080")
  println("main ==========")
}

func initDB() *gorm.DB {
  // 初始化mysql连接
  db, err := gorm.Open(mysql.Open(config.Config.Mysql.DSN), &gorm.Config{})
  if err != nil {
    // panic 会令整个 goroutine 结束
    panic("failed to connect database")
  }
  err = dao.InitTable(db)
  if err != nil {
    panic("failed to init table")
  }

  return db
}

func initServer() *gin.Engine {
  server := gin.Default()

  // cors
  server.Use(cors.New(cors.Config{
    AllowOrigins: []string{"http://localhost:3000"},
    //AllowMethods:     []string{"POST", "GET"},
    // x-token 是前端请求携带的token字段
    AllowHeaders:     []string{"Content-Type", "x-token"},
    ExposeHeaders:    []string{"x-jwt-token"},
    AllowCredentials: true,
    //AllowOriginFunc: func(origin string) bool {
    //  return origin == "http://localhost:3000"
    //},
    MaxAge: 12 * time.Hour,
  }))

  return server
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
  server.Use(sessions.Sessions("wei_session", store))

  // 自定义中间件
  //server.Use(middleware.NewLoginMiddleBuilder().Build())
  server.Use(middleware.NewLoginJWTMiddleBuilder().Build())
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

func initUser(db *gorm.DB) *user.UserHandler {
  dao := dao.NewUserDAO(db)
  repo := repo.NewUserRepo(dao)
  svc := service.NewUserService(repo)
  userHandler := user.NewUserHandler(svc)

  return userHandler
}
