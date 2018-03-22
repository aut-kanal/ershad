FROM alpine:3.5

RUN apk --no-cache add --update ca-certificates
RUN update-ca-certificates

EXPOSE 8000

ADD ershadd /bin/ershadd
WORKDIR /bin

ENTRYPOINT ["/bin/ershadd"]
