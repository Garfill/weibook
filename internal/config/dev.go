//go:build !k8s

package config

var Config = DBConfig{
  Mysql: mysql{
    DSN: "localhost:3306",
  },
  Redis: redis{
    Addr: "localhost:6379",
  },
}
