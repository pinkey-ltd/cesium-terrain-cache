FROM golang:1.24.1-alpine3.21 as build

#ENV GOPROXY
ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn,direct
#build
WORKDIR /go/cache
COPY go.mod .
COPY go.sum .
RUN go mod download
WORKDIR /go/release
COPY . .

RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix cgo -o service ./cmd/pcmp2

FROM alpine:3.19.1 as prod

COPY --from=build /go/release/localtime /etc/localtime
COPY --from=build /go/release/config.yaml /
COPY --from=build /go/release/service /
RUN  mkdir -p /upload/ && \
     mkdir -p /dist/
# ENV configure
EXPOSE 8000
EXPOSE 8080
CMD ["./service"]