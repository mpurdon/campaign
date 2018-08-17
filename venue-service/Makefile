proto:
	protoc -I. --go_out=plugins=micro:${HOME}/Projects/funnelcrowd/gomicro/venue-service proto/venue/venue.proto

build:
	docker build --tag venue-service --build-arg SERVICE=venue .

run:
	docker run -p 50052:50051 -e build_target=venue-service -e MICRO_SERVER_ADDRESS=:50051 -e MICRO_REGISTRY=mdns venue-service
