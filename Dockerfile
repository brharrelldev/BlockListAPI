FROM golang:1.15-alpine

ARG PORT
ARG CACHE_PATH
ARG BLOCKLIST_USER
ARG BLOCKLIST_PASS
ARG DB_PATH

ENV PORT ${PORT}
ENV CACHE_PATH ${CACHE_PATH}
ENV BLOCKLIST_USER ${BLOCKLIST_USER}
ENV BLOCKLIST_PASSWORD ${BLOCKLIST_PASS}
ENV DB_PATH ${DB_PATH}
#
ADD . /app
RUN apk add --update alpine-sdk
RUN cd /app && make build-all


ENTRYPOINT ["/app/bin/blocklist-cli"]