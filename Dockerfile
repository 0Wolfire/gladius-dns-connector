# build stage
FROM golang:1.11.4 AS build-env
ADD . /src
RUN cd /src && go build -tags netgo -a -v -o gladns -ldflags '-w -extldflags "-static"' -i main.go

# final stage
FROM alpine
RUN apk --no-cache --update upgrade && apk --no-cache add ca-certificates
COPY --from=build-env /src/gladns /bin/
