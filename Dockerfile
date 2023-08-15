FROM golang:1.20 AS builder

WORKDIR /app

ENV AWS_SECERT_NAME=clients/creds/bets-ui
ENV AWS_REGION=us-east-1

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -v -o main ./cmd/fundrick/

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .

EXPOSE 8080

# Command to run the executable
CMD ["./main"]
