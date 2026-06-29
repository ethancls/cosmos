# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /cosmos-server ./cmd/server/

# Run stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /cosmos-server /usr/local/bin/cosmos-server
COPY migrations/ /migrations/
EXPOSE 8080
ENTRYPOINT ["cosmos-server"]
