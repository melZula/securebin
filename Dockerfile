FROM golang

RUN mkdir -p /app
WORKDIR /app
RUN cd /app

COPY securebin securebin
COPY 3952.ttf 3952.ttf
COPY configs/apiserver.toml configs/apiserver.toml

ENV CGO_ENABLED=0

EXPOSE 8080

ENTRYPOINT [ "./securebin" ]