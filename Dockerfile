##
## Build
##
FROM golang:1.18-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY *.go ./

RUN CGO_ENABLED=0 go build -o /ec2snapshot

##
## Deploy
##
FROM alpine/curl
WORKDIR /

COPY --from=build /ec2snapshot /ec2snapshot

USER nobody:nobody

ENTRYPOINT ["/ec2snapshot"]