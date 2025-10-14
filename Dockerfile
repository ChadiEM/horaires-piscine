FROM golang:1.25.3-alpine AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

# Build the application
# CGO_ENABLED=0 is important for static linking when using alpine base image
# -ldflags "-s -w" reduces the binary size by stripping debug information
RUN CGO_ENABLED=0 go build -o horaires-piscine -ldflags "-s -w" ./cmd/horaires-piscine

FROM alpine:3.22.1

WORKDIR /

COPY --from=build /app/horaires-piscine .

USER nobody

ENTRYPOINT [ "/horaires-piscine" ]
