FROM jakubknejzlik/godog as builder
WORKDIR /go/src/github.com/jakubknejzlik/godog-sql
COPY . .
RUN apk add --update build-base && go get -t ./...
RUN godog -o /tmp/godog

FROM alpine
VOLUME [ "/godog/features" ]
WORKDIR /godog
COPY --from=builder /tmp/godog /usr/local/bin
CMD [ "godog" ]
