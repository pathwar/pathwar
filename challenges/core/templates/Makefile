CHILDREN = nginx apache-php

all:	build

build release clean:
	for child in $(CHILDREN); do $(MAKE) -C $$child $@; done
