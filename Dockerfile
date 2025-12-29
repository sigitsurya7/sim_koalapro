FROM golang:1.21-alpine AS build

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/seed ./cmd/seed

FROM gcr.io/distroless/static:nonroot

ENV GIN_MODE=release

COPY --from=build /out/api /api
COPY --from=build /out/seed /seed
COPY db/schema.sql /db/schema.sql

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/api"]
