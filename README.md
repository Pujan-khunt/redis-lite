# Redis Lite

A simple TCP server which listens on a configurable port for connections.
It parses RESP style messages and executes command like `SET`, `GET` or `DEL` accordingly.

It strictly understands and parses RESP arrays.
E.g. `*3\r\n$3\r\nSET\r\n$5\r\mykey\r\n$7\r\nmyvalue\r\n`
Which translates to:
- `*3\r\n`: 3 elements in RESP array
- `$3\r\n`: Bulk string of length 3
- `SET\r\n`: Content of bulk string with length exactly 3 (command `SET`)
- `$5\r\n`: Bulk String of length 5

It will update an in-memory database (map for now 😅) and provide an appropriate response.

| Configurable Option | Default Value | Description |
| :--- | :--- | :--- |
| port | 6379 (to honour real redis) | port that the server will bind on |
| ip | 0.0.0.0 | IP address to bind to (0.0.0.0 for all interfaces) |
