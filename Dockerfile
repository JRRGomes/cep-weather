FROM golang:1.20-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
RUN CGO_ENABLED=0 GOOS=linux go build -o /cep-weather-api ./cmd

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /cep-weather-api /cep-weather-api
EXPOSE 8080

CMD ["/cep-weather-api"]
