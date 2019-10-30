FROM alpine:3.10
RUN apk add --no-cache ca-certificates
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
COPY ./firestore_export_debug /firestore_export_debug
ENTRYPOINT ["/firestore_export_debug"]