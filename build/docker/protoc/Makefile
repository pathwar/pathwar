IMAGE ?= pathwar/protoc:v1

.PHONY: release
build:
	docker build --pull -t $(IMAGE) .

.PHONY: release
release: build
	docker push $(IMAGE)
