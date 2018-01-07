FROM golang:1.9

ENV APP="/go/go-api"
WORKDIR $APP

COPY get-dependencies.sh $APP/
RUN chmod +x get-dependencies.sh && ./get-dependencies.sh

ADD src $APP/src
COPY config.yaml $APP/

RUN go build -o bin/go-api src/go-api.go

RUN chmod +x bin/go-api

EXPOSE 8000

ENTRYPOINT ["./bin/go-api"]
