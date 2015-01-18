APIFILE =		apiary.apib
AGLIO_TEMPLATE =	default

dev:
	aglio -i $(APIFILE) -t $(AGLIO_TEMPLATE) -s


release:
	$(MAKE) release_do || $(MAKE) release_teardown


travis:
	find . -name Dockerfile | xargs cat | grep -vi ^maintainer | bash -n
	aglio -i apiary.apib -o apiary.html

release_do:
	git branch -D gh-pages || true
	git checkout -b gh-pages
	aglio -i $(APIFILE) -t $(AGLIO_TEMPLATE) -o index.html
	git add index.html
	git commit index.html -m "Rebuild assets"
	git push -u origin gh-pages -f
	$(MAKE) release_teardown


release_teardown:
	git checkout master
