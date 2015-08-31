NAME ?=		pathwar/unix
VERSIONS ?=	latest onbuild
EXAMPLES ?=	hello
EXTRA_DEPENDENCIES ?=	$(addsuffix /overlay, $(VERSIONS))


-include ../template.mk


$(EXTRA_DEPENDENCIES):	$(shell find ./overlay)
	rm -rf $@
	cp -rf ./overlay $@
