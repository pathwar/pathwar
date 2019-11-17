## moul/rules.mk
INSTALL_STEPS += recursive_install
GENERATE_STEPS += recursive_generate
TEST_STEPS += recursive_test
#BUMPDEPS_STEPS += recursive_bumpdeps
include rules.mk  # see https://github.com/moul/rules.mk
##

.PHONY: recursive_install
recursive_install:
	cd go; make install

.PHONY: recursive_test
recursive_test:
	cd go; make test
	cd web; make test

.PHONY: recursive_bumpdeps
recursive_bumpdeps:
	cd go; make bumpdeps

.PHONY: clean
clean:
	cd go; make clean
	cd docs; make clean

.PHONY: recursive_generate
recursive_generate:
	cd go; make generate
	cd docs; make generate

.PHONY: integration
integration:
	cd tool/integration; make run-suite
