FROM golang:1.12

RUN git clone https://github.com/lotusnoir/powerdns_exporter.git /powerdns_exporter
WORKDIR /powerdns_exporter

RUN go get github.com/prometheus/client_golang/prometheus && \
    go get github.com/sirupsen/logrus && \
    go get gopkg.in/alecthomas/kingpin.v2

RUN make build

ENTRYPOINT ["./powerdns_exporter"]
