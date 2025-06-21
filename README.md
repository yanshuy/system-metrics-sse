Live demo: https://system-metrics.onrender.com/

# System Metrics SSE Streamer

A simple Go application that streams CPU and Memory usage statistics over Server-Sent Events (SSE).

## Running the Application

```bash
go run main.go
```

## Accessing the Stream

The SSE stream is available at: `http://localhost:3333/stream`

You can use tools like `curl` or a web browser with JavaScript to connect:

```bash
curl http://localhost:3333/stream
```
