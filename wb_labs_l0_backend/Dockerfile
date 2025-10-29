FROM golang:1.24-alpine as build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o wb_backend ./cmd/app

FROM alpine:3.19
WORKDIR /app
COPY --from=build /app/wb_backend .
EXPOSE 8080
CMD ["./wb_backend"]