FROM golang:latest AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o sse-server main.go


FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/sse-server .
COPY static/ ./static/
EXPOSE 3333

ENTRYPOINT ["./sse-server"]