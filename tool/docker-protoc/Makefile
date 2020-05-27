IMAGE ?= 	pathwar/protoc
VERSION ?= 	6

build:
	docker build --pull -t $(IMAGE):$(VERSION) .

publish: build
	docker tag $(IMAGE):$(VERSION) $(IMAGE):latest
	docker push $(IMAGE):$(VERSION)
	docker push $(IMAGE):latest
