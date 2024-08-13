FROM 710267309417.dkr.ecr.us-east-1.amazonaws.com/ecr-public/docker/library/alpine:3 AS certs
RUN apk add -U --no-cache ca-certificates make

FROM 710267309417.dkr.ecr.us-east-1.amazonaws.com/ecr-public/docker/library/golang:1.18-alpine AS build
RUN apk add -U --no-cache make
WORKDIR /workspace
COPY . .
RUN make install

FROM 710267309417.dkr.ecr.us-east-1.amazonaws.com/ecr-public/docker/library/alpine:3
EXPOSE 3000
ENV GODEBUG netdns=go
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/bin/drone-ecr-registry /go/bin/

ENTRYPOINT ["/go/bin/drone-ecr-registry"]
