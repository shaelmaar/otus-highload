FROM golang:1.24-alpine AS base
WORKDIR /app

RUN mkdir /out

COPY go.mod .
COPY go.sum .
RUN go mod download

FROM base AS build

WORKDIR /app
ADD . /app

RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    go build \
    -tags go_tarantool_ssl_disable \
    -o /out/social-network-dialogs cmd/*.go

FROM golang:1.24-alpine AS release

WORKDIR /app

COPY --from=build /out/social-network-dialogs /app/

ENTRYPOINT ["/app/social-network-dialogs"]
CMD ["serve"]
