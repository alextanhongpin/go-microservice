FROM golang:1.12-stretch as builder

WORKDIR /go/src/github.com/alextanhongpin/go-microservice

ENV GO111MODULE=on

COPY go.mod go.sum ./

RUN go mod download

COPY ./ .

# If you hit the following error:
#     standard_init_linux.go:190: exec user process caused "no such file or directory"
# It means you did not set CGO_ENABLED=0.
RUN CGO_ENABLED=0 GOOS=linux go build -i -a -installsuffix cgo -o app .

RUN adduser -S -D -H -h /go/src/github.com/alextanhongpin/go-microservice user 
USER user 

FROM alpine:3.9
RUN apk --no-cache add ca-certificates

WORKDIR /app/

COPY --from=builder /go/src/github.com/alextanhongpin/go-microservice/app .

# Allow only tmp folder to be written to. Specify which folder you want to
# allow the application to write to by changing the folder name.
RUN mkdir /app/tmp
RUN adduser -S -D -H -h ./tmp user 
USER user 

# Metadata params
ARG BUILD_DATE
ARG CMD
ARG DESCRIPTION
ARG NAME
ARG URL 
ARG VCS_REF
ARG VCS_URL
ARG VENDOR
ARG VERSION

# Metadata
LABEL org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.name=$NAME \
      org.label-schema.description=$DESCRIPTION \
      org.label-schema.url=$URL \
      org.label-schema.vcs-url=$VCS_URL \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vendor=$VENDOR \
      org.label-schema.version=$VERSION \
      org.label-schema.docker.schema-version="1.0" \
      org.label-schema.docker.cmd=$CMD

CMD ["./app"]
