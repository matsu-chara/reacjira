FROM golang:1.16.0-alpine as builder
MAINTAINER matsu-chara <matsuy00@gmail.com>

WORKDIR /opt/reacjira

COPY . .
RUN apk add --no-cache ca-certificates git
RUN go mod download
RUN go build

FROM alpine
ENV REACJIRA_CONFIG_NAME "/opt/reacjira/config.toml"
ENV REACJIRA_REACJIRA_NAME "/opt/reacjira/reacjira.toml"

RUN apk add --no-cache ca-certificates
COPY --from=builder /opt/reacjira/reacjira /opt/reacjira/reacjira
COPY --from=builder /opt/reacjira/config.toml /opt/reacjira/config.toml
COPY --from=builder /opt/reacjira/reacjira.toml /opt/reacjira/reacjira.toml
ENTRYPOINT ["/opt/reacjira/reacjira"]
