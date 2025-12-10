package config

type mysql struct {
  DSN string
}
type redis struct {
  Addr string
}
type DBConfig struct {
  Mysql mysql
  Redis redis
}
