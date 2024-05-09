FROM golang:1.22-alpine3.19 as build

# 设置国内镜像
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && \
    go env -w GOPROXY=https://goproxy.cn,direct

RUN apk add --no-cache nodejs npm yarn make

WORKDIR /app

COPY . .

RUN cd frontend && yarn install && yarn build

RUN go mod download && go build -o bin/dns-kit cmd/app/main.go

FROM alpine:3.19

WORKDIR /app

COPY --from=build /app/bin /app/dns-kit

CMD ["/app/dns-kit"]
