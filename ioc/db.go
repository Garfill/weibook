package ioc

import (
  "weibook/internal/config"
  "weibook/internal/repo/dao"

  "gorm.io/driver/mysql"
  "gorm.io/gorm"
)

func InitDB() *gorm.DB {
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
