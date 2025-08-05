########
# base #
########

# なぜかbuster以外だと、WASMビルドで真っ白表示になってしまう
FROM golang:1.24.5-bullseye AS base
RUN apt update
RUN apt install -y \
    gcc \
    libc6-dev \
    libgl1-mesa-dev \
    libxcursor-dev \
    libxi-dev \
    libxinerama-dev \
    libxrandr-dev \
    libxxf86vm-dev \
    libasound2-dev \
    pkg-config \
    xorg-dev \
    libx11-dev \
    libopenal-dev \
    upx-ucl

###########
# builder #
###########

FROM base AS builder

WORKDIR /build
COPY go.mod ./
COPY go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GO111MODULE=on \
    go build -o ./bin/ruins .
RUN upx-ucl --best --ultra-brute ./bin/ruins

###########
# release #
###########

FROM gcr.io/distroless/base-debian12:latest AS release

COPY --from=builder /build/bin/ruins /bin/
WORKDIR /work
ENTRYPOINT ["ruins"]

########
# node #
########

FROM node:24.5.0 AS releaser
RUN yarn install
