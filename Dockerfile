FROM golang:1.24.2-alpine as builder

WORKDIR /opt/reacjira

COPY . .
RUN go mod download
RUN go build -o reacjira

FROM alpine:3.21.3
ENV REACJIRA_CONFIG_NAME "/opt/reacjira/config.toml"
ENV REACJIRA_REACJIRA_NAME "/opt/reacjira/reacjira.toml"

COPY --from=builder /opt/reacjira/reacjira /opt/reacjira/reacjira
COPY --from=builder /opt/reacjira/config.toml /opt/reacjira/config.toml
COPY --from=builder /opt/reacjira/reacjira.toml /opt/reacjira/reacjira.toml
RUN apk add --no-cache ca-certificates
ENTRYPOINT ["/opt/reacjira/reacjira"]
