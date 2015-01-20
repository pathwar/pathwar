CHILDREN = nginx apache-php ruby python nodejs golang


all:	build


build release clean examples dependencies:
	for child in $(CHILDREN); do $(MAKE) -C $$child $@; done


dist:	dependencies
	git add -f */*/skeleton
