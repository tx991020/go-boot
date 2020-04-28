FROM alpine:3.4

RUN apk update && apk add --no-cache ca-certificates && \
    apk add tzdata && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

ADD ./{{AppName}} /go/bin/{{AppName}}
WORKDIR /go/src/{{AppName}}
CMD ["/go/bin/{{AppName}}"]
EXPOSE {{http_port}}
