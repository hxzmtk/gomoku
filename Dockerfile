FROM golang:alpine as builder
ENV GO111MODULE on
ENV GOPROXY https://goproxy.io
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
ENV ADDR ":8000"
CMD ["./gomoku"]