build:
	docker build --tag campaign-cli --build-arg SERVICE=campaign .

run:
	docker run -e build_target=campaign-cli -e MICRO_REGISTRY=mdns campaign-cli
