FROM golang:1.8.5 as build
WORKDIR /go/src/github.com/BurtonR/sqlrest

COPY vendor         vendor
COPY handlers       handlers
COPY database       database
COPY middleware     middleware
COPY main.go        .

RUN go get
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o sqlrest .

FROM alpine:3.6
WORKDIR /root/

EXPOSE 5050

COPY --from=build /go/src/github.com/BurtonR/sqlrest/sqlrest    .

CMD ["./sqlrest"]