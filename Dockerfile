FROM golang:1.22-alpine3.19 as build

ARG APP_NAME=dns-kit
ARG APP_VERSION=1.0.0
ARG CHINA_MIRROR=false
ARG LD_FLAGS="-w -s -X github.com/codfrm/cago/configs.Version=${APP_VERSION}"

# 设置国内镜像
RUN if [ "$CHINA_MIRROR" = "true" ]; then \
    echo "Using China Mirror" && \
    sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && \
    go env -w GOPROXY=https://goproxy.cn,direct; \
    fi

RUN apk add --no-cache nodejs npm yarn make

WORKDIR /app

COPY . .

RUN cd frontend && yarn install && yarn build

RUN go mod download && \
    go build -o "bin/${APP_NAME}" -trimpath -ldflags "${LD_FLAGS}" ./cmd/app

FROM alpine:3.19

WORKDIR /app

COPY --from=build /app/bin/dns-kit /app/dns-kit

COPY ./configs/config.yaml.example /app/configs/config.yaml.example

COPY ./deploy/entrypoint.sh /app/entrypoint.sh

RUN chmod +x /app/entrypoint.sh && chmod +x /app/dns-kit

ENTRYPOINT ["/app/entrypoint.sh"]

CMD ["/app/dns-kit","--config","./runtime/config.yaml"]
