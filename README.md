目录结构

```
weibook/
|-- script/ # 脚本（mysql初始化）
|-- internal/
| |-- www/ # api服务
| |-- service/ # 领域服务
| |-- repo/ # 数据存储
| |-- domain/ # 各个领域
```

# 项目打包启动
项目分为k8s环境（线上使用）和docker环境（本地开发）

## k8s环境
1. 运行命令
```shell
GOOS=linux GOARCH=arm go build -tags=k8s  -o weibook .
```
2. 通过 `docker build -t garfill/weibook:v0.0.1 .` 打包镜像
3. **k8s** 开头的 .yaml 文件是用来启动 k8s deployment/service的
   - 依次使用 `kubectl apply -f <file>.yaml` 启动k8s服务

## docker（本地开发）
1. docker-compose.yaml 启动`mysql`和`redis`
2. 运行命令 `go run main.go`