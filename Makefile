.PHONY: mock
mock:
	@mockgen -source=internal/service/user.go -package=svcmocks -destination=internal/service/mock/user.mock.go
	@mockgen -source=internal/service/code.go -package=svcmocks -destination=internal/service/mock/code.mock.go
	@mockgen -source=internal/repo/user.go -package=repomocks -destination=internal/repo/mock/user.mock.go
	@mockgen -source=internal/repo/code.go -package=repomocks -destination=internal/repo/mock/code.mock.go
	@mockgen -source=internal/repo/dao/user.go -package=daomocks -destination=internal/repo/dao/mock/user.mock.go
	@mockgen -source=internal/repo/cache/user.go -package=cachemocks -destination=internal/repo/cache/mock/user.mock.go
	@mockgen -package=redismocks -destination=internal/repo/cache/redismocks/cmd.mock.go github.com/redis/go-redis/v9 Cmdable
	@go mod tidy