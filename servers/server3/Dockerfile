FROM golang:1.21-alpine
WORKDIR /app
COPY main.go .
RUN go mod init server3
RUN go build -o server3 .
CMD ["./server3"]