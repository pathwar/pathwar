NAME ?=		pathwar/sshd
VERSIONS ?=	latest onbuild
EXAMPLES ?=	hello
EXTRA_DEPENDENCIES ?=	$(addsuffix /overlay, $(VERSIONS))


-include ../template.mk


$(EXTRA_DEPENDENCIES):	./overlay
	rm -rf $@
	cp -rf $< $@
