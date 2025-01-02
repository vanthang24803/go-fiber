# Stage 1: Build
FROM golang:1.20-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o app .

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/app .

EXPOSE 4005

CMD ["./app"]
