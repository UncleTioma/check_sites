FROM golang:1.17-alpine as builder

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh make

ENV SERVICE_NAME check_sites
ENV APP /src/${SERVICE_NAME}/
ENV WORKDIR ${GOPATH}${APP}

WORKDIR $WORKDIR

ADD . $WORKDIR

RUN go get ./...
RUN go get -u golang.org/x/lint/golint

RUN go mod tidy
RUN CGO_ENABLED=0 go build -i -v -o release/check_sites

FROM sleeck/crond

COPY --from=builder /go/src/check_sites/release/check_sites /

CMD ["/check_sites"]