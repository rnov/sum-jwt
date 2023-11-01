FROM golang:1 as builder

ENV APP=service

COPY . /${APP}/
WORKDIR /${APP}

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=readonly -a -o /go-rnov/${APP} ./cmd/${APP}

FROM scratch
WORKDIR /

COPY --from=builder /go-rnov/* /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

CMD ["/service"]
