# Build stage
FROM golang:1.23.3-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o backend .

# Final stage
FROM alpine:latest


RUN apk add --no-cache tzdata ca-certificates

ENV TZ=Asia/Bangkok
RUN cp /usr/share/zoneinfo/Asia/Bangkok /etc/localtime && \
    echo "Asia/Bangkok" > /etc/timezone

RUN adduser -D appuser

WORKDIR /app

COPY --from=builder /app/backend .

USER appuser

EXPOSE 8080

ENTRYPOINT ["./backend"]
