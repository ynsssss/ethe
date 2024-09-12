FROM golang:1.23

#hardcoded api key
ENV INFURA_API_KEY=93b9edf4f8364481b183aa34907f1117

RUN mkdir /app
ADD . /app

WORKDIR /app

ENV CGO_ENABLED=1
RUN go env -w GOCACHE=/go-cache
RUN go env -w GOMODCACHE=/gomod-cache
RUN --mount=type=cache,target=/gomod-cache go mod download
RUN --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache \
  go build -o server

EXPOSE 8080

CMD [ "/app/server" ]
