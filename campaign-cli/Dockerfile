#
# Build Container
#

FROM golang:alpine3.8 as go_builder

ARG SERVICE

ENV service $SERVICE
ENV build_target "${service}-cli"

# Need git to pull from github
RUN apk add git && apk add ca-certificates
#RUN apk add protobuf

# Create gouser
RUN adduser -D -g '' gouser

# Using vgo instead of go get, to get around a few
# quirks in how go get works with sub-packages.
#RUN go get -u golang.org/x/vgo

# Create a dep project, and run `ensure`, which will pull in all
# of the dependencies within this directory.
#RUN vgo version
#RUN vgo build
#RUN vgo install ./...
#RUN vgo mod init
#RUN vgo mod sync

# Requirements
RUN go get "github.com/golang/protobuf/proto"
#RUN go get "github.com/golang/protobuf/protoc-gen-go"
RUN go get "github.com/micro/go-micro"
RUN go get "golang.org/x/net/context"
RUN go get "github.com/mpurdon/gomicro-example/campaign-service/proto/campaign"

# Generate the protobuf libs
#RUN protoc -I. --go_out=plugins=micro:. proto/${service}/${service}.proto

RUN echo "Building the binary '${build_target}'"

WORKDIR $GOPATH/src/${build_target}

# Copy the current code into our workdir
COPY . .

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o ${build_target} .

# Strip debugging symbols and stuff from the binary
RUN apk --no-cache add binutils
RUN strip ${build_target}

# Compress the binary
RUN apk --no-cache add upx
RUN upx ${build_target}
RUN chown gouser ${build_target}

#
# Execution Container
#

FROM scratch

ARG SERVICE

ENV service $SERVICE
ENV build_target "${service}-cli"

# Get the CA certs
COPY --from=go_builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the user we created in the build container
COPY --from=go_builder /etc/passwd /etc/passwd

# Here, instead of copying the binary from our host machine,
# we pull the binary from the container named `builder`, within
# this build context. This reaches into our previous image, finds
# the binary we built, and pulls it into this container. Amazing!
COPY --from=go_builder /go/src/${build_target}/${build_target} .
COPY --from=go_builder /go/src/${build_target}/campaign.json .

USER gouser

# Run the binary with a binary build in a separate container,
# with all of the correct dependencies and run time libraries.
#CMD ["./${build_target}"]
ENTRYPOINT ["./campaign-cli"]
#ENTRYPOINT ["./${build_target}"]
