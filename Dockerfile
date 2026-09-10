# build stage
# NOTE: bullseye is archived (upx package gone, old Go can't build go1.27 modules) -> trixie
FROM golang:1.27-trixie as builder

ENV CGO_ENABLED=0

# NOTE: Debian trixie ships UPX 4.x only as `upx-ucl`; the binary is /usr/bin/upx-ucl (symlinked to upx)
RUN apt-get -qq update && \
    apt-get install -yqq upx-ucl && \
    ln -s "$(command -v upx-ucl)" /usr/local/bin/upx

COPY . /build
WORKDIR /build

ARG BUILD_TIME=unknown
ARG GITHASH=unknown
ARG BUILD_TAG=dev

RUN go build -ldflags "-X main.buildStamp=${BUILD_TIME} -X main.gitHash=${GITHASH} -X main.buildTag=${BUILD_TAG}"
RUN strip /build/img-receipt
RUN upx -q -9 /build/img-receipt

# ---
FROM scratch

#COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/img-receipt .

EXPOSE 8080

ENTRYPOINT ["./img-receipt"]
