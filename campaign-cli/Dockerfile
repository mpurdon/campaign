FROM alpine:latest

RUN apk add --update bash

RUN mkdir -p /app
WORKDIR /app

ADD campaign.json /app/campaign.json
ADD campaign-cli /app/campaign-cli

CMD ["./campaign-cli"]
