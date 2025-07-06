FROM golang:1.24-alpine AS base
WORKDIR /app

RUN mkdir /out

COPY go.mod .
COPY go.sum .
RUN go mod download

FROM base AS build

WORKDIR /app
ADD . /app

COPY postgresql/migrations /out/migrations

RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    go build \
    -o /out/social-network cmd/*.go

FROM golang:1.24-alpine AS release

COPY --from=build /out/migrations /app/postgresql/migrations

WORKDIR /app

COPY --from=build /out/social-network /app/

ENTRYPOINT ["/app/social-network"]
CMD ["serve"]
