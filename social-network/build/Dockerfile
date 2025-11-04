FROM golang:1.24-alpine AS base
WORKDIR /app

RUN mkdir /out

COPY go.mod .
COPY go.sum .
RUN go mod download

FROM base AS build

WORKDIR /app
ADD . /app

COPY postgresql/migrations /out/postgresql/migrations
COPY mongo/migrations /out/mongo/migrations

RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    go build \
    -tags go_tarantool_ssl_disable \
    -o /out/social-network cmd/*.go

FROM golang:1.24-alpine AS release

COPY --from=build /out/postgresql/migrations /app/postgresql/migrations
COPY --from=build /out/mongo/migrations /app/mongo/migrations
COPY --from=build /app/resources /app/resources

WORKDIR /app

COPY --from=build /out/social-network /app/

ENTRYPOINT ["/app/social-network"]
CMD ["serve"]
