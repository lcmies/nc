FROM golang:1.20

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

ENV PATH=$PATH:/go/bin
CMD ["nc"]