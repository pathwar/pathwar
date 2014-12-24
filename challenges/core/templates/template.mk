all:	build

example:	1.7.8-onbuild/build
	#$(MAKE) -C example build
	#$(MAKE) -C example run
	$(MAKE) -C example up

build:	$(addsuffix /build, $(VERSIONS))

$(addsuffix /skeleton, $(VERSIONS)):	../../skeleton ../../skeleton/* ../../skeleton/*/*
	../../skeleton/install_local.sh $@

$(addsuffix /build, $(VERSIONS)):	$(@:/build=/skeleton)
	$(eval VERSION := $(shell dirname $@))
	$(eval TAGS := $(shell cat $(VERSION)/tags))
	docker build -t $(NAME):$(VERSION) $(VERSION)
	for tag in $(TAGS); do \
		docker tag $(NAME):$(VERSION) $(NAME):$$tag; \
	done

release:
	docker push $(NAME)
