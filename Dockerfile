FROM alpine

ARG VERSION
ENV UPSTREAM github.com/UncleTioma/check_sites

ENV GOROOT /usr/lib/go
ENV GOPATH /gopath
ENV GOBIN /gopath/bin
ENV PATH $PATH:$GOROOT/bin:$GOPATH/bin

# Install dependencies for building httpdiff 
RUN apk --no-cache update && apk --no-cache upgrade && \
 apk --no-cache add ca-certificates && \
 apk --no-cache add --virtual build-dependencies curl git go musl-dev && \
 # Install check_sites client
 echo "Starting installing check_sites $VERSION." && \
 go get -d $UPSTREAM && \
 cd $GOPATH/src/$UPSTREAM/ && git checkout $VERSION && \
 go install $UPSTREAM && \
 apk del build-dependencies

ENTRYPOINT ["check_sites"]