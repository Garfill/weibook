//go:build k8s

package config

var Config = DBConfig{
  Mysql: mysql{
    DSN: "root:12345678@tcp(weibook-mysql:3306)/weibook?charset=utf8&parseTime=True&loc=Local",
  },
  Redis: redis{
    Addr: "weibook-redis:6379",
  },
}
