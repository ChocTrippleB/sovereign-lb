# sovereign-lb

A TCP load balancer built from scratch in Go. No libraries, no frameworks just the standard library and a mutex.

I built this after reading about and watching a video by Vasilios Syrakis, a Senior Systems Engineer at Atlassian who built Sovereign: an Envoy proxy management system running across 2,000 servers in 13 AWS regions. That kind of infrastructure is the kind of thing I want to understand from the ground up. This is step one.

---

## What it actually does

When a request hits port 8080, the balancer picks the next healthy backend using round-robin and proxies the request there. If a server is down, it skips it. If all servers are down, it returns a 503. When a dead server comes back, the balancer notices within 5 seconds and starts routing to it again no restart needed.

Three backends. One balancer. One goroutine running health checks in the background. That's the whole thing.

```
Client → loadbalancer:8080 → server1:8081
                           → server2:8082
                           → server3:8083
```

---

## How to run it

### With Docker

You need Docker. That's all.

```bash
docker compose up
```

The balancer comes up on `localhost:8080`. The three backend servers start behind it. No Go installation required everything runs in containers.

### Without Docker

You need Go installed. Open four terminal windows.

**Terminals 1–3 start the backend servers:**

```bash
# Terminal 1
go run servers/server1/main.go

# Terminal 2
go run servers/server2/main.go

# Terminal 3
go run servers/server3/main.go
```

**Terminal 4 start the balancer:**

The balancer's backend URLs are hardcoded to Docker hostnames (`server1`, `server2`, `server3`). Swap them to `localhost` in `main.go` before running:

```go
var servers = []*Server{
    {URL: "http://localhost:8081", Alive: true},
    {URL: "http://localhost:8082", Alive: true},
    {URL: "http://localhost:8083", Alive: true},
}
```

Then:

```bash
go run main.go
```

Balancer runs on `localhost:8080`. Same behavior as Docker round-robin, health checks, automatic failover.

---

## Test it

```bash
curl http://localhost:8080
curl http://localhost:8080
curl http://localhost:8080
```

You'll see the response rotate:

```
Hello from Server 1
Hello from Server 2
Hello from Server 3
```

Round-robin. Each request goes to the next server in line.

---

## Kill a server and watch what happens

Stop one of the backend containers from Docker Desktop or with:

```bash
docker stop server2
```

Within 5 seconds the health checker marks it dead. The next request skips it. Traffic keeps flowing between the two remaining servers.

Bring it back:

```bash
docker start server2
```

The balancer detects it's alive again. It re-enters the rotation. Nothing needs to restart.

---

## How the load balancer works internally

**Round-robin selection** `getNextServer()` tracks a `current` index. Every call increments it, wrapping around with modulo. If a server is marked dead, it skips it and tries the next one. If nothing is alive, it returns `nil`.

**Concurrency safety** Multiple requests hit the balancer at the same time. The `current` index is shared state, so it's protected by a `sync.Mutex`. One goroutine at a time gets to read and update it.

**Health checks** A goroutine runs `healthCheck()` in the background from the moment the process starts. Every 5 seconds it pings each backend with an HTTP GET. If the response is anything other than 200 or if the connection fails it flips `server.Alive` to `false`. It logs the transition both ways so you can see servers going down and coming back up in the terminal.

**Reverse proxying** - `httputil.NewSingleHostReverseProxy` does the actual forwarding. It takes a parsed URL, rewrites the incoming request to point at the backend, and handles the response. The balancer is transparent to the client.

---

## Project layout

```
sovereign-lb/
├── main.go                  # Load balancer all the logic lives here
├── Dockerfile               # Container for the balancer
├── docker-compose.yml       # Wires everything together
└── servers/
    ├── server1/
    │   ├── main.go          # "Hello from Server 1"
    │   └── Dockerfile
    ├── server2/
    │   ├── main.go          # "Hello from Server 2"
    │   └── Dockerfile
    └── server3/
        ├── main.go          # "Hello from Server 3"
        └── Dockerfile
```

---

## Stack

- Go (standard library only)
- Docker
- Docker Compose

---

## What I learned building this

Mutexes aren't complicated they just mean "one at a time." The tricky part is remembering *why* you need one: multiple goroutines sharing the same integer is a data race, and data races in Go are silent bugs that corrupt state in ways that are hard to debug.

`httputil.ReverseProxy` does a lot of heavy lifting. Understanding what it's doing under the hood rewriting Host headers, stripping hop-by-hop headers, handling response streaming is the next thing worth digging into.

Health checking is simple but the placement matters. Locking inside the loop, not outside, means the check for each server is atomic but the balancer isn't frozen the whole time.

---

## What's next

- Weighted round-robin give some backends more traffic than others
- Least-connections route to whoever has the fewest active requests
- Sticky sessions pin a client to the same backend based on IP or cookie
- Proper graceful shutdown drain in-flight requests before exiting
- Metrics endpoint expose request counts, error rates, backend health over HTTP
