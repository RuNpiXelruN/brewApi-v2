FROM golang:1.9

ARG app_dev

ENV APP_DEV $app_dev

WORKDIR /go/src/go_apps/go_api_apps/brewApi-v2

ADD . .

RUN go install
RUN go get -u github.com/tools/godep
RUN godep restore
RUN godep go build
RUN go get -u github.com/pilu/fresh

CMD fresh -c fresh.conf;