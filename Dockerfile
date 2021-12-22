FROM golang
LABEL Author ToolmanP
EXPOSE 8080
EXPOSE 80
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
ENV REDIS_ADDR 172.17.0.1:6379
ENV MONGO_ADDR mongodb://172.17.0.1:27017/
WORKDIR /go/src/JiaoNiBan-data
COPY . .
RUN go get github.com/gocolly/colly
RUN go get go.mongodb.org/mongo-driver
RUN go get github.com/go-redis/redis/v8
RUN go get github.com/rabbitmq/amqp091-go
RUN go mod tidy
RUN go build -o dean ./main.go