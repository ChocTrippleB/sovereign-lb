FROM golang:1.21-alpine
WORKDIR /app
COPY main.go .
RUN go mod init loadbalancer
RUN go build -o loadbalancer .
CMD ["./loadbalancer"]