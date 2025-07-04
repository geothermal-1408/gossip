# Gossip - Simple TCP Chat Server

A lightweight TCP chat server implementation in Go that supports basic messaging and rate limiting.

## Features

- TCP-based chat server
- Rate limiting to prevent spam (1 message per second)
- Auto-banning system for repeated rate limit violations
- Simple client disconnect handling
- Basic command support (`:quit`)

## Technical Details

- Port: 6969 (default)
- Rate limit: 1 message per second
- Ban duration: 10 minutes
- Message buffer size: 64 bytes

## Usage

1. Start the server:

```bash
go run main.go
```

2. Connect using a TCP client (e.g., netcat):

```bash
nc localhost 6969
```

3. Type messages to chat
4. Use `:quit` to disconnect

## Rate Limiting

- Users can send 1 message per second
- Exceeding rate limit 3 times results in a 10-minute ban
- Banned IPs are automatically unbanned after the ban duration

## Error Handling

The server includes basic error handling for:

- Connection failures
- Message transmission errors
- Client disconnections
