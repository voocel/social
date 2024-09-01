FROM golang:1.23-alpine AS build

ENV GOPATH=/gobuild
COPY ./app ./app
COPY ./core ./core
COPY ./*.go go.mod LICENSE README.md ./container/app.conf ./

ENV GOPROXY=https://goproxy.cn,direct
RUN rm -rf go.sum && sed -i 's/\/*replace github.com\/ixre/\/\/replace  github.com\/ixre/g' go.mod && \
    go mod tidy && \
    CGO_ENABLED=0 GOOS=linux ARCH=amd64 go build -o social main.go && \
    mkdir -p /opt/social/dist && \
    cp -r social LICENSE README.md config.yaml /opt/social/dist

FROM alpine
MAINTAINER voocel
LABEL vendor="Social"
LABEL version="1.0.0"

WORKDIR /app
COPY --from=build /opt/social/dist/* /app/
COPY ./assets ./assets

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories  \
	&& apk --update add tzdata ca-certificates \
	&& cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
	&& apk del tzdata \
	&& ln -s /app/social /bin \
    && echo "if [ ! -f '/data/config.yaml' ];then cp -r /app/config.yaml /data;fi;" \
    "social -conf /data/config.yaml -d"> /docker-boot.sh
VOLUME ["/data"]
EXPOSE 8800 8100
CMD ["sh","/docker-boot.sh"]