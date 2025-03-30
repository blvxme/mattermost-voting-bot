FROM golang:latest AS builder
WORKDIR /votingbot
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o votingbot cmd/votingbot/main.go

FROM alpine:latest
WORKDIR /votingbot
COPY --from=builder /votingbot/votingbot .
CMD ["./votingbot"]
