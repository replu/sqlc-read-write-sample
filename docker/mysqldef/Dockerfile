FROM debian:bookworm-slim

RUN --mount=type=cache,target=/var/lib/apt,sharing=locked \
    --mount=type=cache,target=/var/cache/apt,sharing=locked \
    apt-get update && apt-get install -y wget

RUN mkdir /sqldef
WORKDIR /sqldef

RUN wget https://github.com/sqldef/sqldef/releases/download/v0.17.20/mysqldef_linux_arm64.tar.gz
RUN tar zxvf mysqldef_linux_arm64.tar.gz

RUN mv mysqldef /usr/local/bin
