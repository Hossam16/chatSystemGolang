
FROM golang:1.13
WORKDIR /
RUN mkdir /go/src/github.com
RUN mkdir /go/src/github.com/Hossam16
RUN mkdir /go/src/github.com/Hossam16/go-chat-creation-api

COPY . /go/src/github.com/Hossam16/go-chat-creation-api
COPY entrypoint.sh /usr/bin/entrypoint-go.sh
RUN chmod +x /usr/bin/entrypoint-go.sh
COPY wait-for-it.sh /usr/bin
RUN chmod +x /usr/bin/wait-for-it.sh
RUN go get -u github.com/go-redis/redis
RUN go get -u github.com/bsm/redislock
RUN go get -u github.com/gorilla/mux 
ENTRYPOINT ["entrypoint-go.sh"]
EXPOSE 8080

CMD go run /go/src/github.com/Hossam16/go-chat-creation-api/cmd/go-chat-creation-api/main.go