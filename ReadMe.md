sovereign-lb/
├── servers/
│   ├── server1/
│   │   └── main.go
│   ├── server2/
│   │   └── main.go
│   └── server3/
│       └── main.go
└── main.go  (balancer - empty for now)


go run servers/server1/main.go

go run servers/server2/main.go

go run servers/server3/main.go
