proto:
	protoc -I. --go_out=plugins=micro:${HOME}/Projects/funnelcrowd/gomicro/account-service proto/account/account.proto

build:
	docker build --tag account-service --build-arg SERVICE=account .

run:
	docker run -p 50053:50051 -e build_target=account-service -e MICRO_SERVER_ADDRESS=:50051 -e MICRO_REGISTRY=mdns account-service
