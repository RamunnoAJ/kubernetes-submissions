FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/random .

FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/random /app/random
USER 65532:65532
CMD ["/app/random"]
