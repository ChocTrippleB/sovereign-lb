FROM golang:1.21-alpine
WORKDIR /app
COPY main.go .
RUN go mod init server2
RUN go build -o server2 .
CMD ["./server2"]