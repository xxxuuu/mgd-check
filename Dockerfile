FROM golang AS builder
WORKDIR /go/src/mgd-check
ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOPROXY https://goproxy.cn,direct
COPY . .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 go build -ldflags="-w -s -extldflags -static" -o /output/server ./main.go

FROM scratch
ENV GIN_MODE release
ENV TZ Asia/Shanghai
COPY --from=builder /output/server /mgd-check/server
WORKDIR /mgd-check/web/template/
COPY --from=builder /go/src/mgd-check/web/template/ .
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
WORKDIR /mgd-check/
EXPOSE 8080
ENTRYPOINT ["./server"]