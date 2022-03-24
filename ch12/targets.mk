generate-ssh-key: ## Generate ssh keys if don't exist
	@(if [ ! -f lab-state/id_rsa ]; then ssh-keygen -b 2048 -t rsa -f ./lab-state/id_rsa -q -N ""; fi)

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