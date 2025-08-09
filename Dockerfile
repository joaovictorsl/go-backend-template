FROM golang:1.25-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the Go application
# CGO_ENABLED=0 disables CGO, creating a statically linked binary
# GOOS=linux specifies the target OS
# -o myapp specifies the output executable name
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main cmd/web/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=build /app/main .

EXPOSE 8000

CMD ["./main"]
