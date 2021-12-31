FROM golang:1.16

ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /data

ADD . /data

#RUN echo cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

RUN go build -o wayne

EXPOSE 1002

CMD [ "/data/wayne" ]
