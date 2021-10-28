ARG GITHUB_PATH=github.com/Damon-V79/act-transition-api

FROM golang:1.17-alpine AS builder

RUN apk add --update make
COPY . /home/${GITHUB_PATH}
WORKDIR /home/${GITHUB_PATH}
RUN make all


FROM alpine:latest AS act-transition-api

LABEL org.opencontainers.image.source https://${GITHUB_PATH}
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /home/${GITHUB_PATH}/bin/act-transition-api .

RUN chown root:root act-transition-api
CMD ["./act-transition-api"]
