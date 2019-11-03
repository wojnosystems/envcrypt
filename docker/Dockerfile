FROM golang:1.13-alpine3.10 AS BUILDER
WORKDIR /tmp/build
COPY go.mod .
COPY go.sum .
RUN go get -d -v all

COPY cmd cmd
COPY pkg pkg
RUN cd /tmp/build/cmd/cli && go build -i -o envcrypt .

FROM alpine:3.10
COPY --from=BUILDER /tmp/build/cmd/cli/envcrypt .
ENTRYPOINT ["envcrypt", "-h"]