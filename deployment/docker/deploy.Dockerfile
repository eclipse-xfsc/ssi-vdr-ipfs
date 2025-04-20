FROM golang:1.21-alpine3.18 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /ipfs

FROM alpine:3.18 as runner
WORKDIR /opt
COPY --from=builder /ipfs .
RUN /ipfs
