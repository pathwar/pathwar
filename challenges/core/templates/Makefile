CHILDREN = nginx apache-php ruby python nodejs golang unix sshd supervisord phantomjs


all:	build


build release clean examples dependencies:
	for child in $(CHILDREN); do $(MAKE) -C $$child $@; done


dist:	dependencies
	git add -f */*/skeleton
	git add -f */*/overlay
