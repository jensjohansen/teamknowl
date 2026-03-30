# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./
# No go.sum yet, we can generate it or just copy if it exists later
# COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /teamknowl-api main.go

# Run stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /teamknowl-api .

EXPOSE 8080

CMD ["./teamknowl-api"]
