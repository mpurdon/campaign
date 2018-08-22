proto:
	protoc -I. --go_out=plugins=micro:${HOME}/Projects/funnelcrowd/gomicro/user-service proto/user/user.proto

build:
	docker build --tag user-service --build-arg SERVICE=user .

run:
	docker run -p 50053:50051 -e build_target=user-service -e MICRO_SERVER_ADDRESS=:50051 -e MICRO_REGISTRY=mdns user-service
