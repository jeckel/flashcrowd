.PHONY: build-image up

build-image:
	@docker-compose build

run:
	@docker-compose up