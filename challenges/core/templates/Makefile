CHILDREN = nginx apache-php

build release:
    for child in $(CHILDREN); do $(MAKE) -C $$child $@; done
