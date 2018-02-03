FROM golang:1.9.3-alpine3.7
MAINTAINER Geoff Lorne <geofflorne@uvic.ca>

RUN apk update && apk upgrade && apk add --no-cache bash git openssh 
RUN mkdir -p /go/src/github.com/moonshot-trading/audit-server
ADD .  /go/src/github.com/moonshot-trading/audit-server
RUN go get github.com/moonshot-trading/audit-server
RUN go install github.com/moonshot-trading/audit-server
ENTRYPOINT /go/bin/audit-server
EXPOSE 44417
