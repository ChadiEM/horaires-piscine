FROM golang:1.25.1-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

ADD *.go ./

RUN go build

FROM alpine:3.22.1

WORKDIR /

COPY --from=build /app/horaires-piscine .

USER nobody

ENTRYPOINT [ "/horaires-piscine" ]
