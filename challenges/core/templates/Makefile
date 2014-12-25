CHILDREN = nginx apache-php

all:	build

build release clean examples:
	for child in $(CHILDREN); do $(MAKE) -C $$child $@; done
