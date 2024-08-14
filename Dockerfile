FROM golang:1.23-alpine as build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

ADD *.go ./

RUN go build

FROM alpine:3.20

WORKDIR /

COPY --from=build /app/horaires-piscine .

USER nobody

ENTRYPOINT [ "/horaires-piscine" ]
