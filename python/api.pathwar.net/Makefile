dev:
	aglio -i api.apib -s


release:
	git branch -D gh-pages || true
	git checkout -b gh-pages
	aglio -i api.apib -t default -o api.html
	git add api.html
	git commit api.html -m "Rebuild assets"
	git push -u origin gh-pages
	git checkout master
