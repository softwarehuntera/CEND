# syntax=docker/dockerfile:1.4
FROM --platform=$BUILDPLATFORM golang:1.23 AS builder

WORKDIR /code

ENV CGO_ENABLED 0
ENV GOPATH /go
ENV GOCACHE /go-build

COPY cend/go.mod cend/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod/cache \
    go mod download

COPY ./cend/ .

RUN --mount=type=cache,target=/go/pkg/mod/cache \
    --mount=type=cache,target=/go-build \
    go build -o bin/cend-backend/cend main.go

CMD ["/code/bin/cend-backend/cend"]

FROM builder as dev-envs

RUN <<EOF
apk update
apk add git
EOF

RUN <<EOF
addgroup -S docker
adduser -S --shell /bin/bash --ingroup docker vscode
EOF

# install Docker tools (cli, buildx, compose)
COPY --from=gloursdocker/docker / /

CMD ["go", "run", "main.go"]

FROM scratch
COPY --from=builder /code/bin/cend-backend/cend /usr/local/bin/cend-backend/cend
CMD ["/usr/local/bin/cend-backend/cend"]