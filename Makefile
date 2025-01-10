.PHONY: plan apply

plan:
	go install .
	cd provider-install-verification && TF_LOG=TRACE TF_LOG_PATH=log.txt terraform plan

apply:
	go install .
	cd provider-install-verification && TF_LOG=TRACE TF_LOG_PATH=log.txt terraform apply
