build:
	docker build --tag user-cli --build-arg SERVICE=user .

run:
	docker run -e build_target=user-cli -e MICRO_REGISTRY=mdns user-cli
