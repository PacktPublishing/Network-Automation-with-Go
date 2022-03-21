tag: check-tag ## Build and tag. Make sure you TAG correctly (Example: export TAG=v0.1.26)
	git add .
	git commit -m "Bump to version ${TAG}"
	git tag -a -m "Bump to version ${TAG}" ${TAG}
	git push --follow-tags

check-tag: ## Check if TAG variable is set.
ifndef TAG
	$(error TAG is undefined)
endif
	@echo "TAG is ${TAG}"

lint: 
	golangci-lint run ./... --disable-all -E errcheck -E lll

# From: https://gist.github.com/klmr/575726c7e05d8780505a
help:
	@echo "$$(tput sgr0)";sed -ne"/^## /{h;s/.*//;:d" -e"H;n;s/^## //;td" -e"s/:.*//;G;s/\\n## /---/;s/\\n/ /g;p;}" ${MAKEFILE_LIST}|awk -F --- -v n=$$(tput cols) -v i=15 -v a="$$(tput setaf 6)" -v z="$$(tput sgr0)" '{printf"%s%*s%s ",a,-i,$$1,z;m=split($$2,w," ");l=n-i;for(j=1;j<=m;j++){l-=length(w[j])+1;if(l<= 0){l=n-i-length(w[j])-1;printf"\n%*s ",-i," ";}printf"%s ",w[j];}printf"\n";}'