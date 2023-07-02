# syntax=docker/dockerfile:1
ARG GO_VERSION=1.20

FROM golang:${GO_VERSION}-alpine as build_go

RUN apk add git 

WORKDIR /app
COPY . /app

ENV GO111MODULE=on
ENV CGO_ENABLED=0

RUN go build -tags urfave_cli_no_docs -ldflags "-X github.com/exler/nurli/internal/cmd.Version=$(git describe --tags)" -o /nurli

FROM alpine:3.18

WORKDIR /app
COPY --from=build_go /nurli /app/nurli
COPY --from=build_go /app/entrypoint.sh /app/entrypoint.sh

EXPOSE 8000

ENTRYPOINT ["./entrypoint.sh"]

CMD ["serve"]
