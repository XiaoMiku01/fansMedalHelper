FROM golang:1.18-alpine

WORKDIR /usr/src/logintool

ARG GOPROXY="https://goproxy.cn"

COPY . .
RUN go get github.com/Baozisoftware/qrcode-terminal-go \
	&& go get github.com/tidwall/gjson \
	&& go build -v -o /usr/local/bin/logintool ./...

CMD ["logintool"]
