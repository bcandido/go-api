FROM golang:1.9

ENV APP="/go/go-api"
WORKDIR $APP

ADD src $APP/src
COPY config.yaml $APP/
COPY get-dependencies.sh $APP/

RUN chmod +x get-dependencies.sh && ./get-dependencies.sh

RUN go build -o bin/go-api src/go-api.go

CMD ['$APP/bin/go-api']
