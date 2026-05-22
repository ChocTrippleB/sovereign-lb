FROM golang:1.21-alpine
WORKDIR /app
COPY main.go .
RUN go mod init server1
RUN go build -o server1 .
CMD ["./server1"]