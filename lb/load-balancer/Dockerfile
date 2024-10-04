FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /lb cmd/lb/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /be cmd/be/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /lb /be ./

EXPOSE 8080 8081

CMD ["./lb"]