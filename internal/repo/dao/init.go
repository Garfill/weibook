package dao

import "gorm.io/gorm"

// 初始化数据表
func InitTable(db *gorm.DB) error {
  // 利用 gorm 的初始化表能力
  // 或者也可以用sql手写建表
  return db.AutoMigrate(&User{})
}
