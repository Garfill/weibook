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

# 项目打包
k8s环境
`
GOOS=linux GOARCH=arm go build -tags=k8s  -o weibook .
`

