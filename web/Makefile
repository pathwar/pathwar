NPM_COMMANDS = test build integration start serve develop format

.PHONY: $(NPM_COMMANDS)
$(NPM_COMMANDS): node_modules
	npm run $@

node_modules: package.json
	npm install
