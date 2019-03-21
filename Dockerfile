FROM alpine:3.6 as alpine
RUN apk add -U --no-cache ca-certificates

FROM alpine:3.6
ENV GODEBUG netdns=go

COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ADD drone-migrate /bin/
ENTRYPOINT [ "/bin/drone-migrate" ]
