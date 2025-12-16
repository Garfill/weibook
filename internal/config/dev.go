//go:build !k8s

package config

var Config = DBConfig{
  Mysql: mysql{
    DSN: "root:12345678@tcp(localhost:3306)/weibook?charset=utf8&parseTime=True&loc=Local",
  },
  Redis: redis{
    Addr: "localhost:6379",
  },
}
