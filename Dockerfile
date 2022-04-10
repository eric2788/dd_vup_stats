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
ENV MYSQL_HOST=192.168.0.127
ENV MYSQL_PORT=3306
ENV MYSQL_USER=ddstats
ENV MYSQL_PASS=changeme
ENV MYSQL_DB=ddstats
ENV WEBSOCKET_URL=ws://192.168.0.127:8888/ws/global
ENV BILIGO_HOST=http://192.168.0.127:8888

EXPOSE 8086

ENTRYPOINT [ "/dd_vup_stats" ]



