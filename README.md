![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![License](https://img.shields.io/badge/license-MIT-blue?style=for-the-badge)

# forgeCache

An in-memory cache server written in Go. Implements the RESP protocol over raw TCP, 
supports string and hash data types, key expiry, and AOF persistence.

## Features

- **RESP protocol** — complete parser and serializer for the Redis wire protocol
- **TCP server** — concurrent client handling via goroutines
- **Commands** — STRING, HASH, and key commands (see full list below)
- **TTL & expiry** — lazy expiry on read + active background cleanup every minute
- **AOF persistence** — write commands logged to disk, replayed on startup to restore state

## Supported Commands

| Command | Description |
|---|---|
| `PING [message]` | Returns PONG or the message |
| `SET key value [EX seconds]` | Set a key with optional expiry |
| `GET key` | Get a string value |
| `DEL key [key ...]` | Delete one or more keys |
| `EXISTS key [key ...]` | Count how many keys exist |
| `INCR key` | Increment integer value |
| `DECR key` | Decrement integer value |
| `APPEND key value` | Append to a string |
| `EXPIRE key seconds` | Set TTL on a key |
| `TTL key` | Get remaining TTL in seconds |
| `HSET key field value [field value ...]` | Set hash fields |
| `HGET key field` | Get a hash field |
| `HDEL key field [field ...]` | Delete hash fields |
| `HGETALL key` | Get all fields and values |
| `HEXISTS key field` | Check if a field exists |
| `HLEN key` | Number of fields in a hash |

## Getting Started

```bash
git clone https://github.com/forge34/forgeCache
cd forgeCache
make build
make run
```

Connect using redis-cli:
```bash
redis-cli -p 4000
```
