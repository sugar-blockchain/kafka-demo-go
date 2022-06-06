ARG GO_VERSION

#FROM golang:${GO_VERSION} as builder
FROM golang:1.16-alpine as builder

WORKDIR /builder/

COPY . .

ENV PATH="/go/bin:${PATH}"
ENV TZ=Asia/Shanghai \
    CGO_ENABLED=1   \
    GO111MODULE=on  \
    GOOS=linux
#    GOPROXY=https://goproxy.cn,direct

# install cgo
RUN apk -U add ca-certificates
RUN apk add build-base
RUN apk add --no-cache bash openssh-client git sudo

# Add ssh private key, to download dependencies from private repositories
RUN mkdir /root/.ssh/
ADD id_rsa.pub /root/.ssh/id_rsa.pub
RUN chmod 400 /root/.ssh/id_rsa.pub

# make sure github domain is accepted
RUN touch /root/.ssh/known_hosts
RUN ssh-keyscan github.com >> /root/.ssh/known_hosts
RUN git config --global --add url."git@github.com:".insteadOf "https://github.com/"
RUN git clone https://github.com/edenhill/librdkafka.git && cd librdkafka && ./configure --prefix /usr && make && make install

#RUN go get ./...
RUN go mod download
RUN go build -tags musl --ldflags "-extldflags -static" -o my_confluent .

FROM scratch
#FROM alpine:latest

ENV TZ=Asia/Shanghai \
    GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
    PROGRAM_ENV=pro

WORKDIR /app

COPY --from=builder /builder/my_confluent .
COPY ./conf ./conf
ADD ca-certificates.crt /etc/ssl/certs/

# for http
EXPOSE 8091

# for rpc
EXPOSE 8092

# 启动服务
CMD ["./my_confluent"]
