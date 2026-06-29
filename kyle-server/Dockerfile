# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /kyle-server ./cmd/server/

# Run stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /kyle-server /usr/local/bin/kyle-server
COPY migrations/ /migrations/
EXPOSE 8080
ENTRYPOINT ["kyle-server"]
