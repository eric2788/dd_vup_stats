FROM golang:1.18-alpine AS builder

WORKDIR /app

COPY . .

RUN apk run update && apk add git

RUN go mod download

RUN go build -o /go/bin/dd_vup_stats

FROM alpine:latest

COPY --from=builder /go/bin/dd_vup_stats /dd_vup_stats
RUN chmod +x /dd_vup_stats

ENV GIN_MODE=release

ENV DB_TYPE=mysql

ENV MYSQL_HOST=172.17.0.1
ENV MYSQL_PORT=3306
ENV MYSQL_USER=ddstats
ENV MYSQL_PASS=changeme
ENV MYSQL_DB=ddstats

ENV PG_HOST=172.17.0.1
ENV PG_PORT=5432
ENV PG_USER=postgres
ENV PG_PASS=changeme
ENV PG_DB=ddstats
ENV PG_SSL=disable

ENV WEBSOCKET_URL=ws://192.168.0.127:8888/ws/global
ENV BILIGO_HOST=http://192.168.0.127:8888
ENV DEV_HOST=http://dev.example.com
ENV REDIS_ADDR=192.168.0.127:6379
ENV REDIS_DB=0

ENV TZ=Asia/Hong_Kong

EXPOSE 8086

ENTRYPOINT [ "/dd_vup_stats" ]



