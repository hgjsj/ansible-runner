FROM golang:1.23.6-alpine3.21 AS builder
RUN mkdir -p /opt/ansible-runner
WORKDIR /opt/ansible-runner
ENV https_proxy="http://proxy-nasa.svcs.entsvcs.com:8080"
ENV http_proxy="http://proxy-nasa.svcs.entsvcs.com:8080"
ENV no_proxy="127.0.0.1,localhost,.env01.mcloud.entsvcs.net"

COPY go.mod go.sum main.go ./
RUN go mod download

COPY internal  ./internal
RUN go build -o ansible-runner 


FROM alpine/ansible:2.17.0
RUN mkdir -p /opt/ansible-runner

COPY --from=builder /opt/ansible-runner/ansible-runner /opt/ansible-runner
COPY cert /opt/ansible-runner/cert
COPY configs /opt/ansible-runner/configs
COPY build_context/ansible-sid-operator /opt/ansible-sid-operator
WORKDIR /opt/ansible-runner
EXPOSE 8080
ENTRYPOINT ["/opt/ansible-runner/ansible-runner"]
CMD ["server"]
