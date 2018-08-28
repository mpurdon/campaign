build:
	docker build --tag user-cli --build-arg SERVICE=user .

run:
	docker run --net="host" -e build_target=user-cli -e MICRO_REGISTRY=mdns user-cli -- --name="Matthew Purdon" --email="mdjpurdon@gmail.com" --password="Testing123"
