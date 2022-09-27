FROM golang:1.19

# Install tdengine client
WORKDIR /
RUN wget https://www.taosdata.com/assets-download/TDengine-client-2.6.0.16-Linux-x64.tar.gz
RUN tar -zxvf /TDengine-client-2.6.0.16-Linux-x64.tar.gz
WORKDIR /TDengine-client-2.6.0.16
RUN chmod +x install_client.sh & ./install_client.sh
RUN rm -rf /TDengine-client-2.6.0.16 & rm -f TDengine-client-2.6.0.16-Linux-x64.tar.gz

WORKDIR /app/
COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w GOOS=linux \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=1 \
    && go env \
    && go mod tidy \
    && go build -o consumer .

CMD ["/app/consumer"]