FROM golang:1.23

RUN mkdir /app
ADD . /app

WORKDIR /app

ENV CGO_ENABLED=1
RUN go env -w GOCACHE=/go-cache
RUN go env -w GOMODCACHE=/gomod-cache
RUN --mount=type=cache,target=/gomod-cache \
  go mod download
RUN --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache \
  go build -o client

CMD [ "/app/client" ]
