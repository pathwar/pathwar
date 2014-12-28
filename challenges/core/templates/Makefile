CHILDREN = nginx apache-php ruby python nodejs

all:	build

build release clean examples:
	for child in $(CHILDREN); do $(MAKE) -C $$child $@; done
