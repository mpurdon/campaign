
build:
	protoc -I. --go_out=plugins=micro:${HOME}/Projects/funnelcrowd/gomicro/campaign-service proto/campaign/campaign.proto
	docker build --tag campaign-service --build-arg SERVICE=campaign .

run:
	docker run -p 50051:50051 \
		-e MICRO_SERVER_ADDRESS=:50051 \
		-e MICRO_REGISTRY=mdns campaign-service
