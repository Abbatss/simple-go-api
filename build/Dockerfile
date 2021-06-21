# Build
FROM golang:1.16-alpine AS build

ARG VERSION=snapshot
ENV VENDORED=true

RUN apk --update --no-cache add make

WORKDIR /app
ADD . /app

RUN make build

# Run
FROM alpine:3.13.4

ENTRYPOINT ["/app/go-app"]
WORKDIR /app
RUN apk --update --no-cache add ca-certificates tzdata && update-ca-certificates

COPY --from=build /app/go-app /app/go-app
