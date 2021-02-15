FROM golang:alpine as builder
ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn
WORKDIR /go/cache
COPY go.mod go.sum ./
RUN go mod download
WORKDIR /app
COPY . .
RUN go build -o gomoku main.go

FROM alpine
WORKDIR /app
COPY web ./web
COPY --from=builder /app/gomoku ./
ENV PORT "8000"
# ENV GIN_MODE "release"
CMD ["sh", "-c", "./gomoku -port ${PORT}"]
