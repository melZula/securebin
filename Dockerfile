# STEP 1

FROM golang:alpine AS builder 

RUN apk update && apk add --no-cache git

ENV USER=appuser
ENV UID=10001 

RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"
WORKDIR $GOPATH/src/melzula/securebin/
COPY . .
RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/securebin -v ./cmd/securebin

# STEP 2

FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /go/bin/securebin /go/bin/securebin

COPY 3952.ttf 3952.ttf
COPY configs/apiserver.toml configs/apiserver.toml

EXPOSE 8080

USER appuser:appuser

ENTRYPOINT ["/go/bin/securebin"]