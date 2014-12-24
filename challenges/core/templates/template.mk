SKELETONS =	$(addsuffix /skeleton, $(VERSIONS))
BUILDS =	$(addsuffix /build, $(VERSIONS))

all:	build

example:	1.7.8-onbuild/build
	#$(MAKE) -C example build
	#$(MAKE) -C example run
	$(MAKE) -C example up

build:	$(BUILDS)

$(BUILDS):	$(SKELETONS)
	$(eval VERSION := $(shell dirname $@))
	$(eval TAGS := $(shell cat $(VERSION)/tags))
	docker build -t $(NAME):$(VERSION) $(VERSION)
	for tag in $(TAGS); do \
		docker tag $(NAME):$(VERSION) $(NAME):$$tag; \
	done

$(SKELETONS):	../../skeleton ../../skeleton/* ../../skeleton/*/*
	../../skeleton/install_local.sh $@

release:
	docker push $(NAME)

clean:
	rm -rf $(SKELETONS)
