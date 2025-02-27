FROM golang:alpine as build
WORKDIR /app
COPY proxy.go .
RUN go build proxy.go

FROM scratch
COPY --from=build /etc/ssl /etc/ssl
COPY --from=build /lib/ld-musl-x86_64.so.1 /lib/ld-musl-x86_64.so.1
COPY --from=build /app/proxy /proxy
USER 1000:1000

CMD ["/proxy"]
