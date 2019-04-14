FROM jakubknejzlik/godog as builder
WORKDIR /go/src/github.com/jakubknejzlik/godog-sql
COPY . .
RUN go get -t ./...
RUN godog -o /tmp/app

FROM alpine
VOLUME [ "/godog/features" ]
WORKDIR /godog
COPY --from=builder /tmp/app /godog/app
CMD [ "/godog/app" ]
