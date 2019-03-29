FROM golang:alpine as build

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh gcc musl-dev

ENV GOROOT=/usr/local/go
COPY . /usr/local/go/src/github.com/mit-dci/bverify-sensor/web
WORKDIR /usr/local/go/src/github.com/mit-dci/bverify-sensor/web
RUN go get -v ./...
RUN go build

FROM node as nodebuild
COPY client/ /tmp/build 
RUN mkdir /tmp/output
WORKDIR /tmp/build
RUN npm install && npm run build
RUN mv build /tmp/output/

FROM golang:alpine
RUN apk add --no-cache ca-certificates
COPY --from=build /usr/local/go/src/github.com/mit-dci/bverify-sensor/web/web /app/bin/web
EXPOSE 8001
WORKDIR /app/bin
COPY --from=nodebuild /tmp/output/build /app/bin/static

ENTRYPOINT ["/app/bin/web"]
