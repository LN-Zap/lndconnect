FROM golang:alpine AS build-env

RUN apk add --update --no-cache git make

RUN git clone --depth 1 https://github.com/LN-Zap/lndconnect.git /lndconnect \
    && cd /lndconnect \
    && go mod download \
    && go get -d -v \
    && make \
    && make install \
    && chmod a+x $GOPATH/bin/lndconnect

FROM alpine

COPY --from=build-env /go/bin/lndconnect /lndconnect

VOLUME [ "/root/.lnd" ]
WORKDIR /

ENTRYPOINT ["/lndconnect"]
