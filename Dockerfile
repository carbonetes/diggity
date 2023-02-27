FROM golang:alpine as build

WORKDIR /diggity

COPY / /diggity

RUN go build .

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

WORKDIR /tmp

COPY --from=build /diggity/diggity /

ENTRYPOINT ["/diggity"]