FROM centos

ENV VERSION 1.8.3
ENV FILE go$VERSION.linux-amd64.tar.gz
ENV URL https://storage.googleapis.com/golang/$FILE
ENV SHA256 0ff3faba02ac83920a65b453785771e75f128fbf9ba4ad1d5e72c044103f9c7a
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN set -eux &&\
  yum -y install git &&\
  yum -y clean all &&\
  curl -OL $URL &&\
	echo "$SHA256  $FILE" | sha256sum -c - &&\
	tar -C /usr/local -xzf $FILE &&\
	rm $FILE &&\
  mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

USER nobody

RUN mkdir -p /go/src/github.com/KamiQuasi/site-reports-api
WORKDIR /go/src/github.com/KamiQuasi/site-reports-api

COPY . /go/src/github.com/KamiQuasi/site-reports-api
RUN go install

ENTRYPOINT /go/bin/site-reports-api

EXPOSE 8080