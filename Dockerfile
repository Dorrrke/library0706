FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o libapp cmd/library0706/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/libapp .
CMD ["./libapp"]
EXPOSE 8080