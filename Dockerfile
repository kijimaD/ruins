########
# base #
########

# なぜかbuster以外だと、WASMビルドで真っ白表示になってしまう
FROM golang:1.20-buster AS base
RUN apt update \
    && apt install -y --no-install-recommends \
    upx-ucl
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
    libopenal-dev

###########
# builder #
###########

FROM base AS builder

WORKDIR /build
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN GO111MODULE=on go build -o ./bin/sokotwo . \
 && upx-ucl --best --ultra-brute ./bin/sokotwo

###########
# release #
###########

FROM gcr.io/distroless/static-debian11:latest AS release

COPY --from=builder /build/bin/sokotwo /bin/
WORKDIR /workdir
ENTRYPOINT ["/bin/sokotwo"]

########
# node #
########

FROM node:18 as releaser
RUN yarn install
