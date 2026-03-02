# Redis Lite

A simple TCP server which listens on a configurable port for connections.
It parses RESP style messages and executes command like `SET`, `GET` or `DEL` accordingly.

Read more about it in [RESP.md](./RESP.md)

It strictly understands and parses RESP arrays.
E.g. `*3\r\n$3\r\nSET\r\n$5\r\mykey\r\n$7\r\nmyvalue\r\n`
Which translates to:
- `*3\r\n`: 3 elements in RESP array
- `$3\r\n`: Bulk string of length 3
- `SET\r\n`: Content of bulk string with length exactly 3 (command `SET`)
- `$5\r\n`: Bulk String of length 5

It will update an in-memory database (map for now 😅) and provide an appropriate response.

## Configuration Options (CLI Flags)

| Configurable Option | Default Value | Description |
| :--- | :--- | :--- |
| port | 6379 (to honour real redis) | port that the server will bind on |
| ip | 0.0.0.0 | IP address to bind to (0.0.0.0 for all interfaces) |

## Supported Redis Commands
| Command | Description | Usage | Expected Response |
| :-- | :-- | :-- | :-- |
| `PING` | check whether connection is alive and or measure latency | `PING` | `PONG` |
| `SET` | create or update key-value pair | `SET "mykey" "myvalue"` | `OK` |
| `GET` | get value from key-value pair | `GET mykey` | `myvalue` |
| `DEL` | delete key-value pair | `DEL mykey` | `1` |

## Usage
Run server using:
```sh
go run .
```

connect using a redis client like `redis-cli` or `valkey-cli`:
```sh
redis-cli
# or valkey-cli
```

Now you can type commands like `SET mykey myvalue`

## Benchmarks

Performance was measured using `redis-benchmark` (or `valkey-benchmark`) running against `redis-lite` on localhost.

**Hardware/Environment:**
* OS: Arch Linux
* CPU: 12th Gen Intel(R) Core(TM) i5-1235U
* RAM: 8GB DDR4-3200

**Test Parameters:**
* Total Requests (`-n`): 100,000
* Concurrent Clients (`-c`): 50
* Payload Size: 3 bytes

| Command | Requests per Second (RPS) | p50 Latency (ms) |
| :--- | :--- | :--- |
| `SET` | ~154000  |  0.143 |
| `GET` | ~189000  |  0.127 |

Run this benchmark locally:
`redis-benchmark -t set,get -q`
