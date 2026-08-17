FROM golang:1.26.6-alpine3.24 AS go-builder
# 走 go mod vendor 模式
ENV CGO_ENABLED=0 GOFLAGS=-mod=vendor
WORKDIR /app
# cache 不常變動的檔案 cmd/ vendor/ go.mod go.sum 加速 build image 用
COPY cmd ./cmd
COPY vendor ./vendor
COPY go.mod .
COPY go.sum .
# 載入業務邏輯及設定
COPY sql ./sql
COPY service ./service
COPY conf.d ./conf.d
ARG buildVersion
ARG buildCommitID
RUN go build -ldflags \
    " \
    -X 'github.com/james730922/wallet/service.BuildVersion=${buildVersion}' \
    -X 'github.com/james730922/wallet/service.BuildCommitID=${buildCommitID}' \
    " \
    -o zqb-apis /app/cmd/zqbapis/

# 只複製執行時所需檔案，降低 image 大小
FROM alpine:3.24.1
RUN apk update && \
    apk upgrade && \
    apk add --no-cache ca-certificates
WORKDIR /zqb
COPY --from=go-builder /app/zqb-apis /zqb/zqb-apis
COPY --from=go-builder /app/conf.d /zqb/conf.d
COPY --from=go-builder /app/sql /zqb/sql
HEALTHCHECK --interval=10s --timeout=3s --start-period=15s --retries=3 \
    CMD wget -q -O /dev/null http://127.0.0.1:17801/readyz || exit 1
ENTRYPOINT ["./zqb-apis"]
