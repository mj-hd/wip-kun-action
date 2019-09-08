FROM golang:1.13.0-alpine as build

WORKDIR /
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s" -a .

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=build /wip-kun /wip-kun

LABEL "com.github.actions.name"="WIP-kun"
LABEL "com.github.actions.description"="manage your pull requests with WIP label"
LABEL "com.github.actions.icon"="shield-off"
LABEL "com.github.actions.color"="yellow"

ENTRYPOINT ["/wip-kun"]
