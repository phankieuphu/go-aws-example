# Use a small base image
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy and build
COPY . .
RUN go mod tidy
RUN go build -o main .

# Final minimal image
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080

CMD ["./main"]
