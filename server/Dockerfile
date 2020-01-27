FROM golang:1.13.6

WORKDIR /go/src/server
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 1323

CMD ["server"]