FROM golang:1.8-alpine

USER nobody

RUN mkdir -p /go/src/github.com/KamiQuasi/site-reports
WORKDIR /go/src/github.com/KamiQuasi/site-reports

COPY . /go/src/github.com/KamiQuasi/site-reports
RUN go-wrapper download && go-wrapper install

CMD ["go-wrapper", "run"]
