FROM golang:alpine

ADD . /go/src/d7024e
ADD ./kademlia /go/bin/kademlia

RUN go install d7024e


WORKDIR /go/bin
ENTRYPOINT d7024e

EXPOSE 8000-8999
EXPOSE 8000-8999/udp
