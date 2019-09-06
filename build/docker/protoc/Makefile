IMAGE ?= pathwar/protoc:v2

.PHONY: release
build:
	docker build --pull -t $(IMAGE) .

.PHONY: release
release: build
	docker push $(IMAGE)
