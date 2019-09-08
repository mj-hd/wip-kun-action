FROM golang:1.13.0-alpine as build

WORKDIR /
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s" -a .

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=build /wip-kun /wip-kun

ENTRYPOINT ["/wip-kun"]