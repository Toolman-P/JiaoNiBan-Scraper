FROM golang
LABEL Author ToolmanP
EXPOSE 8080
EXPOSE 80
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
WORKDIR /app
COPY ./databases ./databases
COPY ./scrapers ./scrapers
COPY ./downloads ./downloads
COPY ./go.sum ./go.sum
COPY ./go.mod ./go.mod
COPY ./main.go ./main.go
RUN go mod tidy
