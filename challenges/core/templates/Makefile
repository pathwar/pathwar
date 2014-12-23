CHILDREN = nginx apache-php

all:	build

build release:
	for child in $(CHILDREN); do $(MAKE) -C $$child $@; done
