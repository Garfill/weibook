FROM ubuntu:latest

# 将编译结果weibook复制到镜像
COPY weibook /app/weibook
WORKDIR /app

ENTRYPOINT ["/app/weibook"]