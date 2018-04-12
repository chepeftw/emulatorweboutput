############################################################
# Build stage
############################################################

FROM golang:1.9 as builder

WORKDIR /go/src/github.com/chepeftw/api/
COPY api.go .
RUN go get -d -v golang.org/x/net/html \
  && go get -d -v gopkg.in/mgo.v2 \
  && go get -d -v github.com/gorilla/mux \
  && go get -d -v github.com/rs/cors
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o myapp .


############################################################
# Run stage
############################################################

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/chepeftw/api/myapp .

COPY wrapper.sh /wrapper.sh
RUN chmod +x /root/myapp
RUN chmod +x /wrapper.sh

CMD /wrapper.sh
