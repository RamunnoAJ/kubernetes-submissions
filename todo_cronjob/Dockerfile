FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/main .
# Install ca-certificates because we are making HTTPS requests to Wikipedia
RUN apk add --no-cache ca-certificates
CMD ["./main"]
