# build stage
FROM golang:1.11.4 AS build-env
ADD . /src
RUN cd /src && go build -tags netgo -a -v -o gladius-dns-connector -ldflags '-w -extldflags "-static"' -i main.go

# final stage
FROM alpine
WORKDIR /app
VOLUME /root/.gladius
COPY --from=build-env /src/gladius-dns-connector /app/
