
FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN GO111MODULE="on" GOOS=linux GOARCH=amd64 go build -o app ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
COPY --from=builder /app/config/config.json config/
COPY --from=builder /app/ui ui
COPY --from=builder /app/forum.db .
EXPOSE 8000
CMD ["./app"]