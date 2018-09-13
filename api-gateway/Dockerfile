FROM alpine:3.2

WORKDIR /

RUN apk add --update ca-certificates && rm -rf /var/cache/apk/* /tmp/*

RUN go get "github.com/micro/go-micro"

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main.go .

# Strip debugging symbols and stuff from the binary
RUN apk --no-cache add binutils
RUN strip ${build_target}

# Compress the binary
RUN apk --no-cache add upx
RUN upx ${build_target}
RUN chown gouser ${build_target}


ENTRYPOINT [ "/micro" ]
