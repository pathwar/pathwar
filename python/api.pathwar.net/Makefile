APIFILE =	pathwar.apib

dev:
	aglio -i $(APIFILE) -s


release:
	git branch -D gh-pages || true
	git checkout -b gh-pages
	aglio -i $(APIFILE) -t default -o index.html
	git add index.html
	git commit index.html -m "Rebuild assets"
	git push -u origin gh-pages -f
	git checkout master
