DEFAULT: help

include .github/targets.mk
include topo-base/targets.mk
include topo-full/targets.mk
include ch07/targets.mk
include ch09/targets.mk
include ch10/targets.mk
include ch12/targets.mk

GOBIN=$(shell which go)


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
cleanup: base-down full-down

## Clone Arista's cEOS image after uploading it
clone:
	docker import cEOS64-lab-4.28.0F.tar ceos:4.28

env-build: generate-ssh-key check-aws-key check-aws-secret ## Build test enviroment on AWS. Make sure you export your API credentials
	@docker run -it \
	--env AWS_ACCESS_KEY_ID \
	--env AWS_SECRET_ACCESS_KEY \
	--volume ${CWD}:/network-automation-with-go:Z \
	ghcr.io/packtpublishing/network-automation-with-go/builder:0.3.3 \
	ansible-playbook /network-automation-with-go/ch12/testbed/create-EC2-testbed.yml \
	--extra-vars "instance_type=$(VM_SIZE) \
	aws_region=$(AWS_REGION) \
	aws_distro=$(AWS_DISTRO)" -v

env-delete: check-aws-key check-aws-secret ## Delete test enviroment on AWS. Make sure you export your API credentials
	@docker run -it \
	--env AWS_ACCESS_KEY_ID \
	--env AWS_SECRET_ACCESS_KEY \
	--volume ${CWD}:/network-automation-with-go:Z \
	ghcr.io/packtpublishing/network-automation-with-go/builder:0.3.3 \
	ansible-playbook /network-automation-with-go/ch12/testbed/delete-EC2-testbed.yml

env-show:  ## Show test environment details
	@cat lab-state/.vm || echo 'VM state file not found'



