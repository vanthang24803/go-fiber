FROM golang:alpine as build

WORKDIR /app

COPY . /app

RUN go build -o app .

FROM alpine

WORKDIR /app

COPY --from=build /app/app .

EXPOSE 4005

CMD ["./app"]
