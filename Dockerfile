# 构建阶段
FROM golang:1.24.2-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /poll-app ./cmd

# 最终镜像
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /poll-app /app/
EXPOSE 8080
CMD ["/app/poll-app"]