FROM golang
 
ADD . /go/src/d7024e
RUN go get github.com/ccding/go-stun
RUN go install d7024e
ENTRYPOINT /go/bin/d7024e
 
EXPOSE 8080

