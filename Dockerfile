FROM golang:alpine3.11 AS builder
MAINTAINER Luc Michalski <michalski.luc@gmail.com>

COPY . /go/src/github.com/lucmichalski/curriculum-vitae-telegram
WORKDIR /go/src/github.com/lucmichalski/curriculum-vitae-telegram

RUN go install

FROM alpine:3.11 AS runtime
MAINTAINER Luc Michalski <michalski.luc@gmail.com>

ARG TINI_VERSION=${TINI_VERSION:-"v0.18.0"}

# Install tini to /usr/local/sbin
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini-muslc-amd64 /usr/local/sbin/tini

# Install runtime dependencies & create runtime user
RUN apk --no-cache --no-progress add ca-certificates \
	&& chmod +x /usr/local/sbin/tini && mkdir -p /opt \
	&& adduser -D lucmichalski -h /opt/feedpushr -s /bin/sh \
	&& su lucmichalski -c 'cd /opt/feedpushr; mkdir -p bin config data'

# Switch to user context
USER lucmichalski
WORKDIR /opt/lucmichalski/data

# copy executable
COPY --from=builder /go/bin/curriculum-vitae-telegram /opt/lucmichalski/bin/curriculum-vitae-telegram

ENV PATH $PATH:/opt/lucmichalski/bin

# Container configuration
# EXPOSE 8080
VOLUME ["/opt/lucmichalski/data"]
ENTRYPOINT ["tini", "-g", "--"]
CMD ["/opt/lucmichalski/bin/curriculum-vitae-telegram"]
