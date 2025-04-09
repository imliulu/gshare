FROM golang:1.23

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app
# 将当前目录同步到docker工作目录下
COPY . .
# RUN CGO_ENABLED=0 GOOS=linux CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOARCH=amd64 go build -o main .
RUN go build .
ENV GIN_MOD=release
EXPOSE 8088

ENTRYPOINT ["./gshare"]
