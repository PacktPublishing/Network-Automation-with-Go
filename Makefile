DEFAULT: lab

include topo-base/targets.mk
include ch07/targets.mk
include ch10/targets.mk

.DEFAULT_GOAL := help

.EXPORT_ALL_VARIABLES:
# AWS Options
AWS_REGION?=us-east-1
AWS_DISTRO?=fedora
VM_SIZE?=t2.large

# Helper variables
CWD=$(shell pwd)
SHELL=/bin/bash

## Cleanup the lab environment
cleanup: 10-down lab-down

## Clone Arista's cEOS image after uploading it
clone:
	docker import cEOS64-lab-4.26.4M.tar ceos:4.26.4M

env-build: generate-ssh-key check-aws-key check-aws-secret ## Build test enviroment on AWS. Make sure you export your API credentials
	@docker run -it \
	--env AWS_ACCESS_KEY_ID \
	--env AWS_SECRET_ACCESS_KEY \
	--volume ${CWD}:/Network-Automation-with-Go \
	ghcr.io/packtpublishing/builder:0.1.25 \
	ansible-playbook create-EC2-testbed.yml \
	--extra-vars "instance_type=$(VM_SIZE) \
	aws_region=$(AWS_REGION) \
	aws_distro=$(AWS_DISTRO)" -v

env-delete: check-aws-key check-aws-secret ## Delete test enviroment on AWS. Make sure you export your API credentials
	@docker run -it \
	--env AWS_ACCESS_KEY_ID \
	--env AWS_SECRET_ACCESS_KEY \
	--volume ${CWD}:/Network-Automation-with-Go \
	ghcr.io/packtpublishing/builder:0.1.25 \
	ansible-playbook delete-EC2-testbed.yml

env-show:  ## Show test environment details
	@cat lab-state/.vm || echo 'VM state file not found'

tag: check-tag ## Build and tag. Make sure you TAG correctly (Example: export TAG=v0.1.26)
	git add .
	git commit -m "Bump to version ${TAG}"
	git tag -a -m "Bump to version ${TAG}" ${TAG}
	git push --follow-tags

generate-ssh-key: ## Generate ssh keys if don't exist
	-ssh-keygen -b 2048 -t rsa -f ./lab-state/id_rsa -q -N "" <<<n 

check-aws-key: ## Check if AWS_ACCESS_KEY_ID variable is set. Brought to you by https://stackoverflow.com/a/4731504
ifndef AWS_ACCESS_KEY_ID
	$(error AWS_ACCESS_KEY_ID is undefined)
endif
	@echo "AWS_ACCESS_KEY_ID is ${AWS_ACCESS_KEY_ID}"

check-aws-secret: ## Check if AWS_SECRET_ACCESS_KEY variable is set.
ifndef AWS_SECRET_ACCESS_KEY
	$(error AWS_SECRET_ACCESS_KEY is undefined)
endif
	@echo "AWS_SECRET_ACCESS_KEY is **************************"

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
