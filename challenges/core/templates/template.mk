SKELETONS =	$(addsuffix /skeleton, $(VERSIONS))
BUILDS =	$(addsuffix /build, $(VERSIONS))
EXAMPLEBUILDS =	$(addprefix examples/, $(addsuffix /build, $(EXAMPLES)))

all:	build

dependencies:	$(SKELETONS)

examples:	$(EXAMPLEBUILDS)

$(EXAMPLEBUILDS):	$(BUILDS)
	$(eval EXAMPLE := $(shell dirname $@))
	cd $(EXAMPLE) && docker-compose build
	cd $(EXAMPLE) && docker-compose up -d

build:	$(BUILDS)

$(BUILDS):	$(SKELETONS)
	$(eval VERSION := $(shell dirname $@))
	$(eval TAGS := $(shell cat $(VERSION)/tags))
	docker build -t $(NAME):$(VERSION) $(VERSION)
	for tag in $(TAGS); do \
		docker tag -f $(NAME):$(VERSION) $(NAME):$$tag; \
	done

$(SKELETONS):	../../skeleton ../../skeleton/* ../../skeleton/*/*
	../../skeleton/install_local.sh $@

release:
	docker push $(NAME)

clean:
	rm -rf $(SKELETONS)
