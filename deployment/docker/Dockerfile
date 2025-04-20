FROM golang:1.21-alpine3.18
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /ipfs
CMD /ipfs
