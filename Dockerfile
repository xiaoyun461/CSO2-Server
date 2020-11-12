#镜像
FROM golang:1.15.3-alpine AS builder
RUN apk add build-base
#安全考虑
#RUN adduser -u 10001 -D app-runner
#设置环境变量
ENV GO111MODULE on
ENV GOPROXY=https://goproxy.cn,direct
#设置工作目录
WORKDIR $GOPATH/src/github.com/KouKouChan/CSO2-Server
#下载mod
COPY go.mod .
COPY go.sum .
RUN go mod download
#复制项目文件
COPY . .
#构建项目
RUN GOOS=linux GOARCH=amd64 go build -o CSO2-Server-docker .
#设置工作目录
WORKDIR $GOPATH/src/github.com/KouKouChan/
#切换可执行文件位置
RUN mv ./CSO2-Server/CSO2-Server-docker ./CSO2-Server-docker
#暴露端口
EXPOSE 1314
EXPOSE 1315
EXPOSE 30001
EXPOSE 30002
#最终运行docker的命令
#USER app-runner
ENTRYPOINT  ["./CSO2-Server-docker"]