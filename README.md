# Redis Lite

A simple TCP server which listens on a configurable port for connections.
Once connected it expects the following 3 types of commands:

1. `SET <key> <value>`
2. `GET <key>`
3. `DEL <key>`

It will update an in-memory database (map for now 😅) and provide an appropriate response.

| Configurable Option | Default Value | Description |
| :--- | :--- | :--- |
| port | 6379 (to honour real redis) | port that the server will bind on |
| ip | 0.0.0.0 | IP address to bind to (0.0.0.0 for all interfaces) |
