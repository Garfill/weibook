.PHONY: mock
mock:
	@mockgen -source=internal/service/user.go -package=svcmocks -destination=internal/service/mock/user.mock.go
	@mockgen -source=internal/service/code.go -package=svcmocks -destination=internal/service/mock/code.mock.go
	@go mod tidy