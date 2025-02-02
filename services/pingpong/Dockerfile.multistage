FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN ls -l /go/pkg/mod/github.com/prometheus/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o pingpong
#final thin image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/pingpong .  
EXPOSE 8080 9080
CMD ["./pingpong"]
