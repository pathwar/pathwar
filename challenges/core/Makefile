level-list:
	curl https://api.github.com/orgs/pathwar/repos 2>/dev/null | jq -r '.[]|select(.name|contains("level"))|["- [", .name, "]", "(http://github.com/pathwar/", .name, ") - ", .description]|join("")'


travis:
	find . -name Dockerfile | xargs cat | grep -vi ^maintainer | bash -n
	find . -name Makefile | xargs make -n
	find . -name "*.mk" | xargs make -n


dist:
	$(MAKE) dist_do || $(MAKE) dist_teardown

dist_do:
	git branch -D dist || true
	git checkout -b dist
	$(MAKE) -C templates/ dist
	git commit -am "Dist"
	git push -u origin dist -f
	$(MAKE) dist_teardown

dist_teardown:
	git checkout master
