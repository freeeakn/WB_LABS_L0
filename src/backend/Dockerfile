FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /backend ./cmd/server

RUN CGO_ENABLED=0 go build -o /publish ./scripts/publish.go

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=builder /backend /backend
COPY --from=builder /publish /publish
COPY model.json /app/model.json
ENTRYPOINT ["/backend"]