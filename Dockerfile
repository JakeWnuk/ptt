FROM golang:alpine as build-env
RUN mkdir src/app
WORKDIR /src/app
ADD ./main.go .
ADD ./go.mod .
ADD ./go.sum .
ADD ./pkg ./pkg
ADD ./templates/ ./templates
RUN go build .

FROM alpine
RUN addgroup --gid 10001 --system nonroot \
    && adduser  --uid 10000 --system --ingroup nonroot --home /home/nonroot nonroot; \
    apk update; apk add --no-cache tini bind-tools

COPY --from=build-env /src/app/ptt /sbin/ptt
ENTRYPOINT ["/sbin/tini", "--", "/sbin/ptt"]
WORKDIR /data
USER nonroot
